package main

import (
	"fmt"
	"crypto/ed25519"
	"context"
	//"io"

	"riddhi/states/p2p"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p/core/network"
	"riddhi/states/protocol/chunk"
	"riddhi/states/protocol/handler"
	"riddhi/states/protocol/packetlayer"
	"riddhi/states/protocol/statemachine"
	"riddhi/states/protocol/transportlayer"
	"riddhi/states/storage"
	"riddhi/states/protocol/session"
	
	
)

type DebugNotify struct{}

func (d *DebugNotify) Listen(network.Network, ma.Multiaddr) {}
func (d *DebugNotify) ListenClose(network.Network, ma.Multiaddr) {}

func (d *DebugNotify) Connected(n network.Network, c network.Conn) {
    fmt.Println("CONNECTED:", c.RemotePeer())
}

func (d *DebugNotify) Disconnected(n network.Network, c network.Conn) {
    fmt.Println("DISCONNECTED:", c.RemotePeer())
}

func (d *DebugNotify) OpenedStream(n network.Network, s network.Stream) {
    fmt.Println("STREAM OPENED:", s.Conn().RemotePeer(), s.Protocol())
}

func (d *DebugNotify) ClosedStream(n network.Network, s network.Stream) {
    fmt.Println("STREAM CLOSED:", s.Conn().RemotePeer(), s.Protocol())
}

func main() {

	ctx := context.Background()

	h, err := p2p.NewHost(
		ctx,
		p2p.HostOptions{
			ListenPort: 4001,
			EnableMDNS: false,
			RelayAddr: "/dns4/cipher-demo.onrender.com/tcp/443/wss/p2p/12D3KooWFvt5h2yNhG7QfLL4CwL6DxsjcaH559Z9eZgBzxgKJqsH",
			PrivKeyPath: "provider.key",
		//	DisableHolePunch:  false, // Enforce strict relay tunneling
		},
	)

	if err != nil {
		panic(err)
	}

	defer h.Close()

	fmt.Println("Provider:", h.ID())

	h.Network().Notify(&DebugNotify{})
	
	// preload chunks
	chunker :=
		chunk.DefaultFileChunker{
			ChunkSize: 32 * 1024,
		}

	chunkedFile, err :=
		chunker.ChunklargeFile(
			"sample.pdf",
		)

	if err != nil {
		panic(err)
	}

	storage.ChunkedFile = chunkedFile
	fmt.Printf( "provider fileID: %x\n", storage.ChunkedFile.FileID, )

	
	//Generate signature
	providerPub, providerPriv, _ := ed25519.GenerateKey(nil)


	p2p.RegisterStreamHandler(
		h,
		func(s network.Stream) {

	fmt.Println("Stream opened from:",
				s.Conn().RemotePeer())
				

	transport := transportlayer.Transport{
	Stream:s,
    }
	defer s.Close()

    providerMachine :=
	statemachine.NewMachine(
		statemachine.Provider,
	)

	providerSession := &session.Session{}
	fmt.Printf(
    "provider session addr: %p\n",
    providerSession,
    )

    //store signature
    providerSession.ProviderPrivateKey = providerPriv
    providerSession.ProviderPublicKey = providerPub


for {

	packet, err := transport.ReceivePacket()

	if err != nil {
		fmt.Println("connection closed:", err)
		return
	}

	// Handle packet based on type
	switch packet.Type {

	case packetlayer.FileMetadataRequest:

	response :=
		handler.HandleFileMetadataRequest(
			packet,
		)

	err =
		transport.SendPacket(
			response,
		)

	if err != nil {
		fmt.Println("send failed")
		return
	}

	case packetlayer.ChunkRequest:

		response :=
			handler.HandleChunkRequest(
				packet,
				providerMachine,
				providerSession,

			)

		err =
			transport.SendPacket(
				response,
			)

		if err != nil {
			fmt.Println(
				"send failed",
			)
			return
		}

	case packetlayer.LotteryTicket:
		

		response :=
			handler.HandleLotteryTicket(
				packet,
				providerMachine,
				providerSession,
			)

		err =
			transport.SendPacket(
				response,
			)

		if err != nil {
			fmt.Println(
				"send failed",
			)
			return
		}
	}
}


		},
	)


	// Keep the provider running forever
	fmt.Println("Waiting for clients...")
	select {}

}