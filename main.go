package main

import (
	"flag"
	"time"

	"github.com/yusufaine/go-tracert/internal/trace"
	"github.com/yusufaine/go-tracert/internal/util"
)

func main() {
	var config util.Config
	flag.StringVar(&config.TargetName, "t", "", "target host, must be supplied")
	flag.IntVar(&config.MaxHops, "h", 32, "max hops, must be greater than 0")
	flag.IntVar(&config.TargetPort, "p", 80, "target port, must be valid and open")
	flag.DurationVar(&config.TimeoutSec, "to", 3*time.Second, "timeout in seconds, must be greater than 0")
	flag.Parse()

	if config.TargetName == "" {
		panic("-t flag required to specify target host")
	} else if config.MaxHops <= 0 {
		panic("-hops flag must be greater than 0")
	} else if config.TimeoutSec <= 0 {
		panic("-timeout flag must be greater than 0")
	}

	config.PopulateConfig()

	trace.Trace(&config)
}
