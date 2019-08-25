# Synopsis

probe each other's aliveness on each side through UDP connection

# Download

    go get github.com/blurty/aliveprobe/connection

# Code Example

    import "github.com/blurty/aliveprobe/connection"

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

    go build

2. modify config.json for your config on each side, and then run like this:

    nohup ./aliveprobe -config config.json &

# Output

```
2019-08-25 10:58:20.606742 +0800 CST m=+5.004018247 probe sequence number: 1 remote machine alive
2019-08-25 10:58:25.609649 +0800 CST m=+10.006942725 probe sequence number: 2 remote machine alive
2019-08-25 10:58:30.609496 +0800 CST m=+15.006807605 probe sequence number: 3 remote machine alive
2019-08-25 10:58:38.614447 +0800 CST m=+23.011787118 probe sequence number: 4 remote machine dead
2019-08-25 10:58:43.614362 +0800 CST m=+28.011718808 probe sequence number: 5 remote machine dead
2019-08-25 10:58:48.610778 +0800 CST m=+33.008151970 probe sequence number: 6 remote machine dead
2019-08-25 10:58:53.605567 +0800 CST m=+38.002958192 probe sequence number: 7 remote machine dead
2019-08-25 10:58:55.609359 +0800 CST m=+40.006757917 probe sequence number: 8 remote machine alive
```