package main

import (
	"fmt"
	"math/rand"
	"time"

	ms "github.com/mitchellh/mapstructure"
)

const (
	losingPingProbability = 0.2
	losingPongProbability = 0
)

type nodeState int

const (
	noToken nodeState = iota
	pingToken
	pongToken
	bothTokens
)

type node struct {
	ping        *token
	pong        *token
	m           int
	state       nodeState
	comm        *comm
	address     string
	nextAddress string
}

func newNode(address string, addresses []string) *node {
	nextAddressIdx := (strIndexOf(address, addresses) + 1) % len(addresses)
	nextAddress := addresses[nextAddressIdx]
	comm := newComm(address, addresses)

	newNode := node{ping: nil, pong: nil, m: 0, state: noToken, comm: comm, address: address, nextAddress: nextAddress}

	if address == strLowestVal(addresses) {
		fmt.Println("Generating first tokens")
		newNode.ping = &token{Type: ping, Value: 1}
		newNode.pong = &token{Type: pong, Value: -1}
	}

	return &newNode
}

func (n *node) run() {

	if n.ping != nil && n.pong != nil {
		fmt.Println("sending first tokens")
		msg := message{Type: tokenMsg, Data: n.ping}
		n.comm.send(msg, n.nextAddress)
		msg = message{Type: tokenMsg, Data: n.pong}
		n.comm.send(msg, n.nextAddress)
	}

	for {
		fmt.Println(n.address, n.state)

		switch n.state {
		case noToken:
			n.listen()
		case pingToken:
			n.criticalSection()
			n.ilisten()
			n.sendPingToken()
		case pongToken:
			n.sendPongToken()
		case bothTokens:
			n.incarnate()
			n.sendPingToken()
			n.sendPongToken()
		}
	}
}

func (n *node) criticalSection() {
	fmt.Println("entering critical section...")
	time.Sleep(time.Second * 2)
	fmt.Println("leaving critical section...")
}

func (n *node) listen() {
	msg := n.comm.recv()
	n.processMsg(msg)
}

func (n *node) ilisten() {
	for {
		msg := n.comm.irecv()
		if msg != nil {
			n.processMsg(msg)
		} else {
			break
		}
	}
}

func (n *node) processMsg(msg *message) {
	switch msg.Type {
	case tokenMsg:
		t := newToken()
		ms.Decode(msg.Data, t)
		fmt.Println("received token", *t)
		n.handleTokenMsg(t)
	default:
		fmt.Println("some other type of message?")
	}
}

func (n *node) handleTokenMsg(t *token) {
	if t.Value == n.m {
		if n.m > 0 {
			fmt.Println("lost Pong token - regenerating", t.Value)
			n.regenerate(t.Value)
			n.state = pongToken
		} else {
			fmt.Println("lost Ping token - regenerating", t.Value)
			n.regenerate(t.Value)
			n.state = pingToken
		}
	} else if abs(t.Value) < abs(n.m) {
		fmt.Println("received some old token?", abs(t.Value), abs(n.m))
		return
	}

	if t.Type == ping {
		n.ping = t

		switch n.state {
		case noToken:
			n.state = pingToken
		case pongToken:
			n.state = bothTokens
		case pingToken, bothTokens:
			fmt.Println("2 ping tokens?")
			panic(nil)
		}
	} else if t.Type == pong {
		n.pong = t

		switch n.state {
		case noToken:
			n.state = pongToken
		case pingToken:
			n.state = bothTokens
		case pongToken, bothTokens:
			fmt.Println("2 pong tokens?")
			panic(nil)
		}
	}
}

func (n *node) regenerate(val int) {
	n.ping = &token{Type: ping, Value: abs(val)}
	n.pong = &token{Type: pong, Value: -n.ping.Value}
}

func (n *node) incarnate() {
	val := abs(n.ping.Value) + 1
	fmt.Println("got 2 tokens - incarnating", val)
	n.ping = &token{Type: ping, Value: val}
	n.pong = &token{Type: pong, Value: -n.ping.Value}
}

func (n *node) sendPingToken() {
	fmt.Println("sending ping token")
	time.Sleep(time.Second)

	if rand.Float64() > losingPingProbability {
		msg := message{Type: tokenMsg, Data: n.ping}
		n.comm.send(msg, n.nextAddress)
	} else {
		fmt.Println("welp, ping token seems to be lost in depths of channel...")
	}

	n.m = n.ping.Value
	n.ping = nil

	switch n.state {
	case pingToken:
		n.state = noToken
	case bothTokens:
		n.state = pongToken
	}
}

func (n *node) sendPongToken() {
	fmt.Println("sending pong token")
	time.Sleep(time.Second * 2)

	if rand.Float64() > losingPongProbability {
		msg := message{Type: tokenMsg, Data: n.pong}
		n.comm.send(msg, n.nextAddress)
	} else {
		fmt.Println("welp, pong token seems to be lost in depths of channel...")
	}

	n.m = n.pong.Value
	n.pong = nil

	switch n.state {
	case pongToken:
		n.state = noToken
	case bothTokens:
		n.state = pingToken
	}
}
