package trace

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/yusufaine/go-tracert/internal/util"
	"golang.org/x/net/icmp" // to de-/serialise ICMP messages, similar to "pack" in python
	"golang.org/x/net/ipv4" // wrapper over syscall for ipv4
)

func Trace(config *util.Config) {
	util.InitMsg(config)

	// Raw socket that listens for ICMP, requires root permissions
	sock, err := icmp.ListenPacket("ip4:icmp", util.IPv4ToString(config.SourceIp))
	if err != nil {
		panic("error creating recv socket: " + err.Error())
	}
	defer sock.Close()

	var (
		id      = os.Getpid() & 0xffff
		reached = false
		seq     = 0
	)
	for ttl := 1; ttl <= config.MaxHops && !reached; ttl++ {
		// Set TTL at every hop
		if err := sock.IPv4PacketConn().SetTTL(ttl); err != nil {
			panic("error setting ttl: " + err.Error())
		}

		// Increment sequence number
		seq++

		// Craft ICMP message, values obtained from wireshark
		icmpBody := &icmp.Echo{
			ID:   id,  // 16 bit ID
			Seq:  seq, // 16 bit sequence number
			Data: make([]byte, 48),
		}
		rand.Read(icmpBody.Data[:]) // fill with random bytes
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho, // 8
			Code: 0,
			Body: icmpBody,
		}

		// Serialise ICMP message
		msgBytes, err := msg.Marshal(nil) // similar to "pack" in python
		if err != nil {
			panic("error serialising icmp message: " + err.Error())
		}

		// Record when message was sent over the wire
		start := time.Now()
		dst := &net.IPAddr{IP: config.TargetIp[:]}
		if _, err := sock.WriteTo(msgBytes, dst); err != nil {
			// Print error and attempt to continue
			fmt.Printf("%d: error detected: %v\n", ttl, err)
			continue
		}

		// Set timeout from now
		if err := sock.SetReadDeadline(time.Now().Add(config.TimeoutSec)); err != nil {
			panic("error setting timeout: " + err.Error())
		}

		// Read reply
		reply := make([]byte, 256)
		n, from, err := sock.ReadFrom(reply)
		if err != nil {
			// Check if error was due to timeout
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				util.PrintOutput(ttl)
			} else { // Unknown error case, try to continue
				util.PrintOutput(ttl, util.WithMsg("unable to read from socket: "+err.Error()))
			}
			continue
		}

		reached = from.String() == util.IPv4ToString(config.TargetIp)

		// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
		replyMsg, err := icmp.ParseMessage(1, reply[:n])
		if err != nil { // Unknown error case, try to continue
			util.PrintOutput(ttl, util.WithMsg("unable to parse ICMP message: "+err.Error()))
			continue
		}

		switch replyMsg.Type {
		case ipv4.ICMPTypeEchoReply, ipv4.ICMPTypeTimeExceeded: // If target reached or TTL hit 0
			echoReply, ok := msg.Body.(*icmp.Echo)
			// Ignore ff reply is not ICMP echo reply, or does not match ID or seq num
			if !ok || echoReply.ID != id || echoReply.Seq != seq {
				continue
			}
			util.PrintOutput(ttl, util.WithHopInfo(from.String(), start))
		}
	}

	if !reached {
		util.PrintOutput(0, util.WithMsg("unable to reach destination, max hops reached"))
	}
}
