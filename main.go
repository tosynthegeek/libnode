package main

import (
	"encoding/json"
	"fmt"
	"io"
	"libp2p/nodes"
	"log"
	"net/http"

	"github.com/libp2p/go-libp2p/core/host"
)
var source host.Host
var target host.Host

type DataPayload struct {
    Data string `json:"data"`
}

func createSourceNode(w http.ResponseWriter, r *http.Request) {
    node, err := nodes.SourceNode()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    source = node
    fmt.Fprintf(w, "Source node created with ID: %s\n", source.ID().String())
}

func createTargetNode(w http.ResponseWriter, r *http.Request) {
    node, err := nodes.TargetNode()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    target = node
    fmt.Fprintf(w, "Target node created with ID: %s\n", target.ID().String())
}

func connectNodes(w http.ResponseWriter, r *http.Request) {
    if source == nil || target == nil {
        http.Error(w, "Source or target node not initialized", http.StatusBadRequest)
        return
    }
    nodes.ConnectNode(source, target)
    fmt.Fprintf(w, "Nodes connected successfully")
}

func getSourceNodePeers(w http.ResponseWriter, r *http.Request) {
    if source == nil {
        http.Error(w, "Source node not initialized", http.StatusBadRequest)
        return
    }
    peers, err := nodes.SourceNodePeers(source)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(peers)
}

// Send data handler
func sendDataHandler(w http.ResponseWriter, r *http.Request) {
    if source == nil || target == nil {
        http.Error(w, "Source or target node not initialized", http.StatusBadRequest)
        return
    }

    var payload DataPayload
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    err = json.Unmarshal(body, &payload)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    err = nodes.SendData(source, target, payload.Data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Datasent: %s\n", payload.Data)

    fmt.Fprintf(w, "Data sent successfully")
}

// Initialize node and receive data handler
func initReceiveHandler(w http.ResponseWriter, r *http.Request) {
    if source == nil && target == nil {
        http.Error(w, "No node initialized", http.StatusBadRequest)
        return
    }

    if source == nil || target == nil{
        fmt.Fprintf(w, "Source node and Target node must be declared and connected\n")
        return
    }
    
    data:= nodes.RecieveData(target)
    if data == "" {
        fmt.Fprint(w, "No data received from source\n")
        return
    }

    fmt.Fprintf(w, "Data Received: %s\n", data)
}

func main() {
    http.HandleFunc("/createSourceNode", createSourceNode)
    http.HandleFunc("/createTargetNode", createTargetNode)
    http.HandleFunc("/connectNodes", connectNodes)
    http.HandleFunc("/getSourceNodePeers", getSourceNodePeers)
	http.HandleFunc("/sendData", sendDataHandler)
	http.HandleFunc("/initReceive", initReceiveHandler)

    fmt.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}