package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	addresses := []string{"127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"}

	nodes := []*node{}
	for idx, addr := range addresses {
		node := newNode(addr, addresses)
		nodes = append(nodes, node)

		if idx == len(addresses)-1 {
			node.run()
		} else {
			go node.run()
		}
	}
}
