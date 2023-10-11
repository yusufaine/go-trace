package trace

import (
	"fmt"
	"syscall"
	"time"

	"github.com/yusufaine/go-tracert/internal/util"
)

func Trace(config Config) {
	targetIp := util.ResolveTargetHostname(config.Target)
	fmt.Printf("Tracing route to %s [%s] from %s:%d\n",
		config.Target,
		util.IPv4ToString(targetIp),
		util.IPv4ToString(config.SourceIp),
		config.SourcePort)

	var found bool
	for ttl := 1; ttl <= config.MaxHops && !found; ttl++ {
		startTime := time.Now()
		// create tcp socket
		sendSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
		if err != nil {
			panic("error creating send socket: " + err.Error())
		}
		defer syscall.Close(sendSock)

		// set ttl
		if err := syscall.SetsockoptInt(sendSock, 0x0, syscall.IP_TTL, ttl); err != nil {
			panic("error setting ttl: " + err.Error())
		}

		// set timeout
		if err := syscall.SetsockoptTimeval(sendSock, syscall.SOL_SOCKET, syscall.SO_SNDTIMEO, &syscall.Timeval{Sec: 1}); err != nil {
			panic("error setting timeout: " + err.Error())
		}

		// create raw icmp socket
		recvSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
		if err != nil {
			panic("error creating recv socket: " + err.Error())
		}
		defer syscall.Close(recvSock)

		// set timeout
		if err := syscall.SetsockoptTimeval(recvSock, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &syscall.Timeval{Sec: 1}); err != nil {
			panic("error setting timeout: " + err.Error())
		}

		// bind to host address
		hostSa := syscall.SockaddrInet4{Addr: config.SourceIp, Port: config.SourcePort}
		if err := syscall.Bind(recvSock, &hostSa); err != nil {
			panic("error binding recv socket: " + err.Error())
		}

		// connect to target
		targetSa := syscall.SockaddrInet4{Addr: targetIp, Port: config.Port}
		if err := syscall.Connect(sendSock, &targetSa); err != nil {
			if err == syscall.ECONNREFUSED {
				util.PrintOutput(ttl, targetIp, 0)
				break
			}

			// receive ICMP packet
			var buf [512]byte
			_, from, err := syscall.Recvfrom(recvSock, buf[:], 0)
			if err != nil {
				panic("error receiving icmp packet: " + err.Error())
			}
			util.PrintOutput(ttl, from.(*syscall.SockaddrInet4).Addr, time.Since(startTime))
			continue
		} else {
			found = true
			util.PrintOutput(ttl, targetIp, time.Since(startTime))
		}
	}

	if !found {
		fmt.Println("Unable to reach target, max hops reached")
	}
}
