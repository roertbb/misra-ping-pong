package main

import (
	"encoding/json"
	"fmt"

	"github.com/pebbe/zmq4"
)

type comm struct {
	recvSock *zmq4.Socket
	socks    map[string]*zmq4.Socket
}

func newComm(address string, addresses []string) *comm {
	newComm := comm{recvSock: nil, socks: map[string]*zmq4.Socket{}}

	recv, _ := zmq4.NewSocket(zmq4.PULL)
	recv.Bind(fmt.Sprintf("tcp://%s", address))
	newComm.recvSock = recv

	for _, addr := range addresses {
		socket, _ := zmq4.NewSocket(zmq4.PUSH)
		newComm.socks[addr] = socket
		newComm.socks[addr].Connect(fmt.Sprintf("tcp://%s", addr))
	}

	return &newComm
}

type msgType int

const (
	tokenMsg msgType = iota
)

type message struct {
	Type msgType     `json:"type"`
	Data interface{} `json:"data"`
}

func (c *comm) send(msg message, address string) {
	mmsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("failed to marshal message")
		panic(err)
	}

	if _, err := c.socks[address].SendBytes([]byte(mmsg), zmq4.DONTWAIT); err != nil {
		fmt.Println("failed to send message")
		panic(err)
	}
}

func (c *comm) recv() *message {
	data, _ := c.recvSock.RecvBytes(0)

	m := message{}
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		fmt.Println("failed to unmarshal message")
		panic(err)
	}

	return &m
}

func (c *comm) irecv() *message {
	data, _ := c.recvSock.RecvBytes(zmq4.DONTWAIT)
	if len(data) == 0 {
		return nil
	}

	m := message{}
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		fmt.Println("failed to unmarshal message")
		panic(err)
	}

	return &m
}
