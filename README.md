# Synopsis

probe each other's aliveness on each side through UDP connection

# Download

    go get github.com/blurtheart/aliveprobe/connection

# Code Example

    import "github.com/blurtheart/aliveprobe/connection"

    // create config from file
    cfg, err := connection.NewConfigFromFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}
    // create connection
	conn, err := connection.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	conn.Start()

# config

    "local_ip":"127.0.0.1", // ip of local side
    "local_port":"12344",   // port of local side
    "remote_ip":"127.0.0.1",    // ip of other side, must be a literal IP address, or a host name that can be resolved to IP addresses
    "remote_port":"12345",      // port of other side, must be a literal port number or a service name
    "period":5,     // duration seconds between two probe datagram
    "timeout":3     // wait seconds before ack of probe

some examples for ("ip:port"):

    ("golang.org:http")
    ("192.0.2.1:http")
    ("198.51.100.1:80")
    ("[2001:db8::1]:domain")
    ("[fe80::1%lo0]:53")
    (":80")

# Usage

1. compile your code

    go build main.go

2. modify config.json for your config on each side, and then run like this:

    ./main -config config.json