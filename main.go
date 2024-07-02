package main

import (
	"libp2p/nodes"
	"log"
)

func main() {
	source, err:= nodes.SourceNode()
	if err != nil {
		log.Fatal(err.Error())
	}
	nodes.NodeID(source)
	nodes.NodeAddress(source)

	target, err:= nodes.TargetNode()
	if err != nil {
		log.Fatal(err.Error())
	}
	nodes.NodeID(target)
	nodes.NodeAddress(target)
}