package nodes

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

func SourceNode() (host.Host, error) {
	node, err:= libp2p.New()
	if err != nil {
		return nil, fmt.Errorf("error creating a new node %w", err)
	}

	return node, nil
}

func TargetNode() (host.Host, error) {
	node, err:= libp2p.New(libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/8007",
	))
	if err != nil {
		return nil, fmt.Errorf("error creating a new node %w", err)
	}

	return node, err
}

func ConnectNode(sourceNode host.Host, targetNode host.Host)  {
	targetNodeAddeInfo := host.InfoFromHost(targetNode)
	err:= sourceNode.Connect(context.Background(), *targetNodeAddeInfo)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func SourceNodePeers(sourceNode host.Host) (int, error) {
	peers:= sourceNode.Network().Peers()
	fmt.Println("Peers: ", peers)

	return len(peers), nil
}

func NodeAddress(node host.Host) {
	addresses:= node.Addrs()
	for _, address:= range addresses {
		fmt.Println("Address: ", address.String()) 
	}
}

func NodeID(node host.Host) {
	fmt.Printf("Node ID: %s", node.ID().String())
}