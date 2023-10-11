# gotrace

> [!NOTE]
> A golang implementation of traceroute using raw sockets which requires root excess.

## Usage

* (Recommended) Download the binary from in the `bin` folder (only built for amd64 darwin (MacOS), linux, and windows), or
* Run it yourself using `go run main.go`, or building it yourself and running the binary (requires golang installed)

```text
Usage of ./bin/gotrace_darwin_amd64:
  -hops int
        max hops, must be greater than 0 (default 32)
  -port int
        target port, must be valid port number (default 80)
  -t string
        target host, must be supplied
  -timeout duration
        timeout in seconds, must be greater than 0 (default 3s)
```

## Acknowledgement

* [pro-bing](https://github.com/prometheus-community/pro-bing/blob/main/ping.go) -- usage of raw sockets and icmp protocol.
