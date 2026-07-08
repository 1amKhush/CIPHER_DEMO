// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"riddhi/states/p2p"
// )

// func main() {

// 	ctx := context.Background()

// 	h, err := p2p.NewHost(
// 		ctx,
// 		p2p.HostOptions{
// 			ListenPort:  4002,
// 			PrivKeyPath: "peer.key",
// 			EnableMDNS:  true,
// 			RelayAddr:   "",
// 		},
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Host started!")

// 	fmt.Println("Peer ID:")
// 	fmt.Println(h.ID())

// 	fmt.Println()

// 	fmt.Println("Listening addresses:")

// 	for _, addr := range h.Addrs() {
// 		fmt.Println(addr)
// 	}
// 	go func() {
// 	for {
// 		time.Sleep(3 * time.Second)

// 		fmt.Println("Connected peers:")

// 		for _, id := range h.Network().Peers() {
// 			fmt.Println(" -", id)
// 		}
// 	}
// }()

// 	select {}
// }
//peers := h.Network().Peers()

	// if len(peers) == 0 {
	// 	panic("no peer found")
	// }
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"riddhi/states/p2p"
	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {
	ctx := context.Background()

	relayCircuitAddr := "/dns4/cipher-demo.onrender.com/tcp/443/wss/p2p/12D3KooWFvt5h2yNhG7QfLL4CwL6DxsjcaH559Z9eZgBzxgKJqsH"

	providerID, err := peer.Decode("12D3KooWNpk5WRASNGsia1XKK8nEVRVWMVN35knKHYRfsxTQ3bL8")
	if err != nil {
		panic(err)
	}

	h, err := p2p.NewHost(
		ctx,
		p2p.HostOptions{
			ListenPort:  4002,
			EnableMDNS:  false,
			RelayAddr:   "/dns4/cipher-demo.onrender.com/tcp/443/wss/p2p/12D3KooWFvt5h2yNhG7QfLL4CwL6DxsjcaH559Z9eZgBzxgKJqsH",
			PrivKeyPath: "client.key",
			//DisableHolePunch:  false, // Enforce strict relay tunneling
		},
	)
	if err != nil {
		panic(err)
	}

	p2p.PrintConnections(h, providerID)

	fmt.Println("Client:", h.ID())

	// Give the host a moment to register on the network
	time.Sleep(5 * time.Second)

	

	// Fixed variable name mismatch here
	connErr := p2p.ConnectViaRelay(
		ctx,
		h,
		relayCircuitAddr,
		providerID,
	)
	if connErr != nil {
		log.Fatalf("Provider failed to connect to relay: %v", connErr)
	}

	fmt.Println("Connectedness:", h.Network().Connectedness(providerID))
	// Start monitoring connections
go func() {
    for {
        fmt.Println("===== Active Connections =====")
        p2p.PrintConnections(h, providerID)
        time.Sleep(2 * time.Second)
    }
}()

// Give libp2p time to perform AutoNAT / DCUtR
time.Sleep(20 * time.Second)

	fmt.Println("Provider Addrs:")
	for _, a := range h.Peerstore().Addrs(providerID) {
		fmt.Println(a)
	}

	fmt.Println("Host Addrs:")
	for _, a := range h.Addrs() {
		fmt.Println(a)
	}

	streamCtx, streamCancel := context.WithTimeout(ctx, 30*time.Second)
	defer streamCancel()

	// Fixed trailing token compilation typo here
	stream, err := p2p.OpenStream(
		streamCtx,
		h,
		providerID,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Operation timed out, retrying or aborting...")
		}
		log.Fatalf("Failed to open stream: %v", err)
	}

	fmt.Println("Remote Addr:")

fmt.Println(
    stream.Conn().RemoteMultiaddr(),
)
	_, err = stream.Write([]byte("Hello Provider"))
	if err != nil {
		log.Printf("Failed to write to stream: %v", err)
		stream.Close()
		return
	}
	// 2. CRITICAL: Close the write-half of the stream.
	// This flushes the buffer and sends an EOF to the provider's Read(),
	// telling the provider "I am done sending data, now I am ready to listen."
	if err := stream.CloseWrite(); err != nil {
		log.Printf("Failed to close write stream: %v", err)
	}

	buf := make([]byte, 1024)
	n, err := stream.Read(buf) // Capture the read error safely
	if err != nil {
		log.Printf("Failed to read from stream: %v", err)
		stream.Close()
		return
	}

	fmt.Println("Reply:", string(buf[:n]))
	stream.Close()
}
