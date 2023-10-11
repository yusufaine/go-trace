package util

import (
	"fmt"
	"net"
	"time"
)

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

func PrintOutput(ttl int, ipAddr [4]byte, timeDiff time.Duration) {
	if len(ipAddr) == 0 {
		fmt.Println(ttl, "*")
	} else {
		ipString := IPv4ToString(ipAddr)
		PadIpAddress(&ipString, 15) // 15 is the length of the longest IPv4 address
		fmt.Printf("%d: %s\t%.2fms\n", ttl, ipString, float32(timeDiff.Microseconds())/1000)
	}
}

// Takes in a string pointer and adds spaces to the end of the string until it is the desired length
func PadIpAddress(str *string, desiredLength int) {
	for len(*str) < desiredLength {
		*str += " "
	}
}
