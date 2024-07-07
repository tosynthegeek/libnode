package nodes

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

func Lnode() {
	node, err:= libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Ping(false),
	)

	if err != nil {
        panic(err)
    }

    // configure our own ping protocol
    pingService := &ping.PingService{Host: node}
    node.SetStreamHandler(ping.ID, pingService.PingHandler)

	peerInfo := peer.AddrInfo{
        ID:    node.ID(),
        Addrs: node.Addrs(),
    }
    addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
    fmt.Println("libp2p node address:", addrs[0])
}

func SourceNode() (host.Host, error) {
	node, err:= libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Ping(false),
	)

	fmt.Println("Node listenig on: ", node.Addrs())
	if err != nil {
		return nil, fmt.Errorf("error creating a new node %w", err)
	}

	// wait for a SIGINT or SIGTERM signal
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
        <-ch
        fmt.Println("Received signal, shutting down...")

        // shut the node down
        if err := node.Close(); err != nil {
                panic(err)
        }

	return node, nil
}

func TargetNode() (host.Host, error) {
	node, err:= libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Ping(false),
	)

	fmt.Println("Node listenig on: ", node.Addrs())
	if err != nil {
		return nil, fmt.Errorf("error creating a new node %w", err)
	}

	// wait for a SIGINT or SIGTERM signal
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
        <-ch
        fmt.Println("Received signal, shutting down...")

        // shut the node down
        if err := node.Close(); err != nil {
                panic(err)
        }

	return node, err
}

func ConnectNode(sourceNode host.Host, targetNode host.Host) error {
	targetNodeAddeInfo := host.InfoFromHost(targetNode)
	err:= sourceNode.Connect(context.Background(), *targetNodeAddeInfo)
	if err != nil {
		return fmt.Errorf("failed to connect nodes: %w", err)
	}
	log.Printf("Source node connected to target node: %s", targetNode.ID().String())
	return nil
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