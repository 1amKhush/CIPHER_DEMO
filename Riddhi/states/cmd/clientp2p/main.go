package main

import (
	"fmt"
	"os"
	"context"
	"time"
	"log"
	"errors"
	"crypto/ed25519"

	"riddhi/states/p2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"riddhi/states/protocol/handler"
	"riddhi/states/protocol/packetlayer"
	"riddhi/states/protocol/statemachine"
	"riddhi/states/protocol/transportlayer"
	"riddhi/states/protocol/session"
)

func main() {

		//connect to provider

		ctx := context.Background()

		relayCircuitAddr := "/dns4/cipher-demo.onrender.com/tcp/443/wss/p2p/12D3KooWFvt5h2yNhG7QfLL4CwL6DxsjcaH559Z9eZgBzxgKJqsH" //Relay address after running relay

		providerID, err := peer.Decode("12D3KooWHriVDgxoZZbX56SGxCUD9aUPn7i1e945Ph2HU3BssCKR") //inside peer.Decode add ProviderID generated after running provider
		if err != nil {
			panic(err)
		}

		h, err := p2p.NewHost(
			ctx,
			p2p.HostOptions{
				ListenPort:  4002,
				EnableMDNS:  false,
				RelayAddr:   "/dns4/cipher-demo.onrender.com/tcp/443/wss/p2p/12D3KooWFvt5h2yNhG7QfLL4CwL6DxsjcaH559Z9eZgBzxgKJqsH",//Relay address after running relay
				PrivKeyPath: "client.key",
				//DisableHolePunch:  false, // Enforce strict relay tunneling
			},
		)
		if err != nil {
			panic(err)
		}

		defer h.Close()

		//debugging
		p2p.PrintConnections(h, providerID)

		fmt.Println("ClientID:", h.ID())

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

		defer stream.Close()

		fmt.Println("Remote Addr:")

	fmt.Println(
		stream.Conn().RemoteMultiaddr(),
	)


	//create a output file to write the rebuilt file
	outputFile, err :=
		os.Create("rebuilt.pdf")

	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	// Create a new state machine for the client
    clientMachine :=
	statemachine.NewMachine(
		statemachine.Requester,
	)

	// Create a new session for the client
	clientSession := &session.Session{}
	fmt.Printf(
    "client session addr: %p\n",
     clientSession,
    )

	//Generate signature
	clientPub, clientPriv, _ := ed25519.GenerateKey(nil)
    //store signature
    clientSession.RequesterPrivateKey = clientPriv
    clientSession.RequesterPublicKey = clientPub

    //request for fileiNfo
	transport := transportlayer.Transport{
    Stream:stream,
     }
	req :=
	packetlayer.FileMetadataRequestPayload{
		FileName: "sample.pdf",
	}

    encoded :=
	packetlayer.EncodeFileMetadataRequest(
		req,
	)

    packet :=
	packetlayer.Packet{
		Type: packetlayer.FileMetadataRequest,
		Length: uint32(len(encoded)),
		Payload: encoded,
	}

	err = transport.SendPacket(&packet)
   if err != nil {
	panic(err)
    }

    responsePacket, err := transport.ReceivePacket()

   if err != nil {
	panic(err)
    }

    metadata :=
	packetlayer.DecodeFileMetadataResponse(
		responsePacket.Payload,
	)

    fmt.Printf(
	"received fileID: %x\n",
	metadata.FileID,
     )

    fmt.Println(
	"file size:",
	metadata.FileSize,
     )

    fmt.Println(
	"chunk size:",
   	metadata.ChunkSize,
     )

	 clientSession.FileID = metadata.FileID


	//nonce is for encryption(not signature)
	var nonce [16]byte
	copy(
		nonce[:],
		[]byte("abcdefghijklmnop"),
	)

	totalChunks := int(
    (metadata.FileSize + uint64(metadata.ChunkSize) - 1) /
    uint64(metadata.ChunkSize),
	)

	clientSession.PaidChunks =
	handler.GeneratePaymentSchedule(
		totalChunks,
	)

	fmt.Println("Payment schedule:")

for chunk := range clientSession.PaidChunks {
	fmt.Println(
		"Paid chunk:",
		chunk,
	)
}

for i := 0; i < totalChunks; i++ {

	req :=
		packetlayer.ChunkRequestPayload{
				ChunkIndex: uint32(i),
				Nonce:      nonce,
		}

	encodedReq :=
			packetlayer.EncodeChunkRequest(
				req,
		)

	requestPacket :=
			packetlayer.Packet{
				Type:
					packetlayer.ChunkRequest,

				Length:
					uint32(len(encodedReq)),

				Payload:
					encodedReq,
		}

	clientMachine.Transition(
			statemachine.WaitingForChunkResponse,
	    )

	fmt.Println(
			"Client requesting chunk",
			i,
		)

    err = transport.SendPacket(&requestPacket)
    responsePacket, err :=
	transport.ReceivePacket()
	if err != nil {
		panic(err)
	}
	if responsePacket == nil {
		panic("nil response packet")
	}

   //handle chunk response
	handler.HandleChunkResponse(
		responsePacket,
		clientMachine,
		clientSession,
	    )

 for{


	// generate ticket
	ticketPacket := handler.CreateLotteryTicket(clientSession)

	transport.SendPacket(ticketPacket)

	revealPacket, err := transport.ReceivePacket()

	plaintext := handler.HandleKeyReveal(
		revealPacket,
		clientMachine,
		clientSession,
	)

	if plaintext == nil {
		fmt.Println(
			"verification failed",
		)
		continue
	}

    _, err =
	outputFile.Write(
		plaintext,
	)

    if err != nil {
	panic(err)
    }

	  break

	

  }


}

	fmt.Printf( "client fileID: %x\n", clientSession.FileID, )

	fmt.Println(
		"file rebuilt",
	)

}