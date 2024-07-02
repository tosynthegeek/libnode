package nodes

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
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

func SendData(sourceNode host.Host, targetNode host.Host, data string) error {
	targetId:= targetNode.ID()
	stream, err:= sourceNode.Network().NewStream(context.Background(), targetId)
	if err != nil {
		return fmt.Errorf("failed to create new stream: %w", err)
	}
    defer stream.Close()

    _, err = stream.Write([]byte(data))
    if err != nil {
        return fmt.Errorf("failed to write to stream: %w", err)
    }

    return nil
}

func RecieveData(node host.Host) chan string {
	dataChan := make(chan string)
    
	node.Network().SetStreamHandler(func(s network.Stream) {
		defer s.Close()

		buf:= make([]byte, 1024)
		n, err:= s.Read(buf)
		if err != nil {
			log.Printf("Error reading from stream: %s", err)
            return
		}

		data:= string(buf[:n])
		log.Printf("Data received oo: %s\n", data)

		dataChan <- data
	})

	return dataChan
}