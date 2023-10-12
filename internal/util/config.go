package util

import (
	"fmt"
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
	d := net.Dialer{Timeout: c.TimeoutSec}
	con, err := d.Dial("tcp", targetIp.String()+":"+fmt.Sprint(c.TargetPort))
	if err != nil {
		panic("unable to connect to target address: " + targetIp.String() + ":" + fmt.Sprint(c.TargetPort))
	}
	defer con.Close()

	// Save in config
	copy(c.SourceIp[:], con.LocalAddr().(*net.TCPAddr).IP.To4())
	c.SourcePort = con.LocalAddr().(*net.TCPAddr).Port
}
