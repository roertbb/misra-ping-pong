package main

import (
	"fmt"
	"math/rand"
	"time"

	ms "github.com/mitchellh/mapstructure"
)

type node struct {
	ping         *token
	pong         *token
	m            int
	comm         *comm
	address      string
	nextAddress  string
	log          *log
	pingLossProb float64
	pongLossProb float64
}

func newNode(address string, addresses []string, pingLossProb, pongLossProb float64) *node {
	idx := strIndexOf(address, addresses)
	nextAddressIdx := (idx + 1) % len(addresses)
	nextAddress := addresses[nextAddressIdx]
	comm := newComm(address, addresses)
	log := newLog(address, idx)

	newNode := node{ping: nil, pong: nil, m: 0, comm: comm, address: address, nextAddress: nextAddress, log: log, pingLossProb: pingLossProb, pongLossProb: pongLossProb}

	if address == strLowestVal(addresses) {
		newNode.log.info("generating first tokens")
		newNode.ping = &token{Type: ping, Value: 1}
		newNode.pong = &token{Type: pong, Value: -1}
	}

	return &newNode
}

func (n *node) run() {

	if n.ping != nil && n.pong != nil {
		n.log.info("sending first tokens")

		n.comm.send(message{Type: tokenMsg, Data: n.ping}, n.nextAddress)
		n.ping = nil
		n.comm.send(message{Type: tokenMsg, Data: n.pong}, n.nextAddress)
		n.pong = nil
	}

	for {
		n.log.debug(fmt.Sprintf("state: ping %s | pong %s ", n.ping.tokenToValue(), n.pong.tokenToValue()))

		if n.ping != nil && n.pong != nil {
			// both token
			n.incarnate()
			n.sendPingToken()
			n.sendPongToken()

		} else if n.ping != nil {
			// ping token
			n.criticalSection()
			n.ilisten()
			n.sendPingToken()

		} else if n.pong != nil {
			// pong token
			n.sendPongToken()

		} else {
			// no token
			n.listen()
		}
	}
}

func (n *node) criticalSection() {
	n.log.info("entering critical section...")
	time.Sleep(time.Second * 2)
	n.log.info("leaving critical section...")
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
		n.log.info(fmt.Sprintf("received %s token with %d value", t.tokenToType(), t.Value))
		n.handleTokenMsg(t)
	default:
		n.log.warn("some other type of message?")
	}
}

func (n *node) handleTokenMsg(t *token) {
	if t.Value == n.m {
		if n.m > 0 {
			n.log.warn("detect PONG token loss - regenerating with value: ", t.Value)
			n.regenerate(t.Value)
		} else {
			n.log.warn("detect PING token loss - regenerating with value: ", t.Value)
			n.regenerate(t.Value)
		}
	} else if abs(t.Value) < abs(n.m) {
		n.log.warn("received some old token?", abs(t.Value), abs(n.m))
		return
	} else if t.Type == ping {
		if n.ping != nil {
			n.log.error("2 PING tokens?")
			panic(nil)
		}

		n.ping = t
	} else if t.Type == pong {
		if n.pong != nil {
			n.log.error("2 PONG tokens?")
			panic(nil)
		}

		n.pong = t
	}
}

func (n *node) regenerate(val int) {
	n.ping = &token{Type: ping, Value: abs(val)}
	n.pong = &token{Type: pong, Value: -n.ping.Value}
}

func (n *node) incarnate() {
	val := abs(n.ping.Value) + 1
	n.log.warn("got 2 tokens - incarnating with value: ", val)
	n.ping = &token{Type: ping, Value: val}
	n.pong = &token{Type: pong, Value: -n.ping.Value}
}

func (n *node) sendPingToken() {
	n.log.info("sending PING token")
	time.Sleep(time.Second)

	if rand.Float64() > n.pingLossProb {
		msg := message{Type: tokenMsg, Data: n.ping}
		n.comm.send(msg, n.nextAddress)
	} else {
		n.log.warn("PING token seems to be lost in depths of channel...")
	}

	n.m = n.ping.Value
	n.ping = nil
}

func (n *node) sendPongToken() {
	n.log.info("sending PONG token")
	time.Sleep(time.Second * 2)

	if rand.Float64() > n.pongLossProb {
		msg := message{Type: tokenMsg, Data: n.pong}
		n.comm.send(msg, n.nextAddress)
	} else {
		n.log.warn("PONG token seems to be lost in depths of channel...")
	}

	n.m = n.pong.Value
	n.pong = nil
}
