package util

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func InitMsg(config *Config) {
	var sb strings.Builder
	sb.WriteString("Tracing route to ")
	sb.WriteString(config.TargetName + ":" + fmt.Sprint(config.TargetPort))
	if config.TargetName != IPv4ToString(config.TargetIp) {
		sb.WriteString(fmt.Sprintf(" [%s:%d]", IPv4ToString(config.TargetIp), config.TargetPort))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("from %s:%d, ", IPv4ToString(config.SourceIp), config.SourcePort))
	sb.WriteString(fmt.Sprintf("max hop %d, ", config.MaxHops))
	sb.WriteString(fmt.Sprintf("timeout %s\n", config.TimeoutSec.String()))
	sb.WriteString("\n")
	fmt.Printf(sb.String())
}

func ResolveTargetHostname(host string) [4]byte {
	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		panic("error resolving ip address: " + err.Error())
	}
	ret := [4]byte{}
	copy(ret[:], ipAddr.IP.To4())
	return ret
}

func IPv4ToString(ipAddr [4]byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", ipAddr[0], ipAddr[1], ipAddr[2], ipAddr[3])
}

type printStruct struct {
	ipAddr    string
	rtt       time.Duration
	customMsg string
}

type printOption func(ps *printStruct)

func WithHopInfo(ipAddr string, start time.Time) printOption {
	return func(ps *printStruct) {
		ps.ipAddr = ipAddr
		ps.rtt = time.Since(start)
	}
}

func WithMsg(msg string) printOption {
	return func(ps *printStruct) {
		ps.customMsg = msg
	}
}

func PrintOutput(ttl int, opts ...printOption) {
	var ps printStruct
	for _, opt := range opts {
		opt(&ps)
	}

	if len(ps.customMsg) > 0 {
		if ttl == 0 {
			fmt.Printf("  %s\n", ps.customMsg)
		} else {
			fmt.Printf("  %d: %s\n", ttl, ps.customMsg)
		}
		return
	}

	if len(ps.ipAddr) == 0 {
		fmt.Printf("  %d: *\n", ttl)
		return
	} else if net.ParseIP(ps.ipAddr).To4() == nil {
		panic(fmt.Sprintf("expected IPv4, received: %v", ps.ipAddr))
	}

	PadIpAddress(&ps.ipAddr, 15) // 15 is the length of the longest IPv4 address
	fmt.Printf("  %d: %s\t%.2fms\n", ttl, ps.ipAddr, float32(ps.rtt.Microseconds())/1000)
}

// Takes in a string pointer and adds spaces to the end of the string until it is the desired length
func PadIpAddress(str *string, desiredLength int) {
	for len(*str) < desiredLength {
		*str += " "
	}
}
