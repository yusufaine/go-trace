package main

import (
	"flag"

	"github.com/yusufaine/go-tracert/internal/trace"
)

func main() {
	var config trace.Config
	flag.StringVar(&config.Target, "t", "1.1.1.1", "target host")
	flag.IntVar(&config.Port, "port", 80, "target port, must be valid port number")
	flag.IntVar(&config.MaxHops, "hops", 15, "max hops, must be greater than 0")
	flag.Parse()

	if config.Target == "" {
		panic("-t flag required to specify target host")
	} else if config.MaxHops <= 0 {
		panic("-hops flag must be greater than 0")
	}

	config.PopulateSourceInfo()

	trace.Trace(config)
}
