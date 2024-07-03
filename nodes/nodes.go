package nodes

import (
	"context"
	"fmt"
	"io"
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

func SendData(source host.Host, target host.Host, data string) error {
	targetId:= target.ID()
	stream, err:= source.Network().NewStream(context.Background(), targetId)
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

func RecieveData(node host.Host) string {
	dataChan := make(chan string, 10)
	var data string
    
	node.Network().SetStreamHandler(func(stream network.Stream) {
		defer stream.Close()

		buf:= make([]byte, 1024)
		for {
            n, err := stream.Read(buf)
            if err != nil {
                if err != io.EOF {
                    log.Printf("Error reading from stream: %s", err)
                }
                close(dataChan) // Close channel when stream ends
                break
            }
			
			log.Println(n)
			data:= string(buf[:n])
			log.Printf("Data received oo: %s\n", data)

			select {
				case dataChan <- data:
					log.Println("Data sent to channel")
				default:
					log.Println("Data channel full or closed")
			}
		}
	})

	return data
}