package trace

import (
	"fmt"
	"syscall"
	"time"

	"github.com/yusufaine/go-tracert/internal/util"
)

func Trace(config Config) {
	targetIp := util.GetIpAddr(config.Target)
	fmt.Printf("Tracing route to %s [%s]\n", config.Target, util.IPv4ToString(targetIp))

	ttl := 0
	for {
		startTime := time.Now()
		if ttl == config.MaxHops {
			panic("Max hops reached")
		}
		ttl++
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

		// Create ICMP socket
		recvSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
		if err != nil {
			panic("error creating recv socket: " + err.Error())
		}
		defer syscall.Close(recvSock)

		// Bind to any address
		if err := syscall.Bind(recvSock, &syscall.SockaddrInet4{}); err != nil {
			panic("error binding recv socket: " + err.Error())
		}

		// Connect to target
		sa := syscall.SockaddrInet4{Port: config.Port, Addr: targetIp}
		if err := syscall.Connect(sendSock, &sa); err != nil {
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
			util.PrintOutput(ttl, targetIp, time.Since(startTime))
			break
		}
	}
}
