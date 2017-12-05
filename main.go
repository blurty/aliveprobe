// Copyright 2017 Blurt. All rights reserved.
// Use of this source code is governed by a GPL-3.0
// license that can be found in the LICENSE file.
//
// alive probe

package main

import (
	"flag"
	"github.com/blurtheart/aliveprobe/connection"
	"log"
)

func main() {
	// read config
	configPath := flag.String("config", "./config/config.json", "path of the config file")
	flag.Parse()

	// create connection config
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
}
