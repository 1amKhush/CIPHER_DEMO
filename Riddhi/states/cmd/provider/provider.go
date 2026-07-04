package main

import (
	"fmt"
	"crypto/ed25519"

	"riddhi/states/protocol/chunk"
	"riddhi/states/protocol/handler"
	"riddhi/states/protocol/packetlayer"
	"riddhi/states/protocol/statemachine"
	"riddhi/states/protocol/transportlayer"
	"riddhi/states/storage"
	"riddhi/states/protocol/session"
)

func main() {

	listener, err :=
		transportlayer.StartServer(
			":8080",
		)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		"provider listening on :8080",
	)


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


	conn, err := listener.Accept()

	if err != nil {
		panic(err)
	}

	transport := transportlayer.Transport{
	Conn: conn,
    }
	defer conn.Close()

    providerMachine :=
	statemachine.NewMachine(
		statemachine.Provider,
	)

	providerSession := &session.Session{}
	fmt.Printf(
    "provider session addr: %p\n",
    providerSession,
    )

	//Generate signature
	providerPub, providerPriv, _ := ed25519.GenerateKey(nil)
    //store signature
    providerSession.ProviderPrivateKey = providerPriv
    providerSession.ProviderPublicKey = providerPub


for {

	packet, err :=
		transport.ReceivePacket()

	if err != nil {
		fmt.Println(
			"connection closed",
		)
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
}