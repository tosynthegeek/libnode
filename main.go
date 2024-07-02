package main

import (
	"fmt"
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

	nodes.ConnectNode(source, target)
	fmt.Println("Nodes connected....")

	peers, err:= nodes.SourceNodePeers(source)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Source node peers: %d", peers)
}