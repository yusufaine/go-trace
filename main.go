package main

import (
	"flag"

	"github.com/yusufaine/go-tracert/internal/trace"
)

func main() {
	var config trace.Config
	flag.StringVar(&config.Target, "t", "", "target host")
	flag.IntVar(&config.Port, "port", 80, "port")
	flag.IntVar(&config.MaxHops, "hops", 15, "max hops")
	flag.Parse()

	if config.Target == "" {
		panic("-t flag required to specify target host")
	}

	trace.Trace(config)
}
