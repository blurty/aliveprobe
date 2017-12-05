// Copyright 2017 Blurt. All rights reserved.
// Use of this source code is governed by a GPL-3.0
// license that can be found in the LICENSE file.
//
// connection between two machines

package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

type Data struct {
	Direction string `json:"direction"` // "0"-上行 "1"-下行
	SN        int    `json:"sn"`        // sequence number
}

var (
	bufferSize = 50
	upstream   = "0"
	downstream = "1"
)

var (
	resolveAddrError = errors.New("failed to parse address")
	dialError        = errors.New("dial udp address failed")
	listenError      = errors.New("listen udp address failed")
)

type Connection struct {
	sConn   *net.UDPConn
	rConn   *net.UDPConn
	rCh     chan int
	period  int
	timeout int
}

type Config struct {
	LocalIP    string `json:"local_ip"`
	RemoteIP   string `json:"remote_ip"`
	LocalPort  string `json:"local_port"`
	RemotePort string `json:"remote_port"`
	Period     int    `json:"period"`
	Timeout    int    `json:"timeout"`
}

func NewConfigFromFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func New(cfg *Config) (*Connection, error) {
	rAddr, err := net.ResolveUDPAddr("udp4", cfg.RemoteIP+":"+cfg.RemotePort)
	if err != nil {
		return nil, resolveAddrError
	}
	rSock, err := net.DialUDP("udp4", nil, rAddr)
	if err != nil {
		return nil, dialError
	}
	// put here avoid socket not reclaimed when rSock error
	lAddr, err := net.ResolveUDPAddr("udp4", cfg.LocalIP+":"+cfg.LocalPort)
	if err != nil {
		return nil, resolveAddrError
	}
	lSock, err := net.ListenUDP("udp4", lAddr)
	if err != nil {
		return nil, listenError
	}
	return &Connection{
		sConn:   rSock,
		rConn:   lSock,
		rCh:     make(chan int, 1),
		period:  cfg.Period,
		timeout: cfg.Timeout,
	}, nil
}

func (conn Connection) probe() error {
	// create send ticker
	tic := time.NewTicker(time.Second * time.Duration(conn.period))
	defer tic.Stop()

	var start int
	data := Data{upstream, start}
	for _ = range tic.C {
		start++
		data.SN = start
		sData, _ := json.Marshal(data)
		conn.sConn.Write(sData)
		tim := time.NewTimer(time.Second * time.Duration(conn.timeout))
	outLoop:
		for {
			select {
			case rSN := <-conn.rCh:
				if rSN != start+1 {
					// unmatched sn
					// fmt.Println(time.Now(), "wrong response, abort data")
				} else {
					// alive
					fmt.Println(time.Now(), "probe sequence number:", start, "remote machine alive")
					break outLoop
				}
			case <-tim.C:
				// dead
				fmt.Println(time.Now(), "probe sequence number:", start, "remote machine dead")
				break outLoop
			}
		}
		tim.Stop()
	}
	return nil
}

func (conn Connection) response(sn int) error {
	data := Data{downstream, sn}
	sData, _ := json.Marshal(data)
	conn.sConn.Write(sData)
	return nil
}

func (conn Connection) recv() {
	sData := make([]byte, bufferSize)
	var data Data
	for {
		dataCount, err := conn.rConn.Read(sData)
		if err != nil || dataCount == 0 {
			// fmt.Println(time.Now(), "read socket error:", err)
			continue
		}
		err = json.Unmarshal(sData[:dataCount], &data)
		if err != nil {
			// fmt.Println(time.Now(), "parse data failed:", err)
			continue
		}
		if data.Direction == upstream {
			conn.response(data.SN + 1)
		} else if data.Direction == downstream {
			conn.rCh <- data.SN
		} else {
			// log here
			// fmt.Println(time.Now(), "data format error:", data)
		}
	}
}

func (conn Connection) Start() {
	go conn.probe()
	conn.recv()
}
