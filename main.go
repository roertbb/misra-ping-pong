package main

import (
	"flag"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	pingLossProb := flag.Float64("ping-loss-prob", 0, "probability of ping token loss in channel")
	pongLossProb := flag.Float64("pong-loss-prob", 0, "probability of pong token loss in channel")

	flag.Parse()

	addresses := flag.Args()

	nodes := []*node{}
	for idx, addr := range addresses {
		node := newNode(addr, addresses, *pingLossProb, *pongLossProb)
		nodes = append(nodes, node)

		if idx == len(addresses)-1 {
			node.run()
		} else {
			go node.run()
		}
	}
}
