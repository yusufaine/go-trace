package util

import (
	"net"
	"time"
)

type Config struct {
	MaxHops    int
	SourceIp   [4]byte
	SourcePort int
	TimeoutSec time.Duration
	TargetIp   [4]byte
	TargetName string
	TargetPort int
}

func (c *Config) PopulateConfig() {
	// Resolve hostname to IP address
	targetIp, err := net.ResolveIPAddr("ip", c.TargetName)
	if err != nil {
		panic("unable to resolve target address: " + err.Error())
	}
	// Save in config
	copy(c.TargetIp[:], targetIp.IP.To4())

	// Let OS choose source IP and port
	con, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: targetIp.IP, Port: 80})
	if err != nil {
		panic("unable to dial target address: " + err.Error())
	}
	defer con.Close()

	// Save in config
	copy(c.SourceIp[:], con.LocalAddr().(*net.TCPAddr).IP.To4())
	c.SourcePort = con.LocalAddr().(*net.TCPAddr).Port
}
