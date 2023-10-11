package trace

import "net"

type Config struct {
	SourceIp   [4]byte
	SourcePort int
	Target     string
	Port       int
	MaxHops    int
}

// Let the OS populate the source IP and port
func (c *Config) PopulateSourceInfo() {
	targetAddr, err := net.ResolveIPAddr("ip", c.Target)
	if err != nil {
		panic("unable to resolve target address: " + err.Error())
	}

	if con, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: targetAddr.IP, Port: 80}); err != nil {
		panic("unable to dial target address: " + err.Error())
	} else {
		defer con.Close()
		copy(c.SourceIp[:], con.LocalAddr().(*net.TCPAddr).IP.To4())
		c.SourcePort = con.LocalAddr().(*net.TCPAddr).Port
	}
}
