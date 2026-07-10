 package main

 /*This code is used to test file transfer using TCP if you want to use it then need to make some changes in protocol:-
In transport and DecodePacket:-instead of Stream network.stream use conn net.Conn
 */

// import (
// 	"fmt"
// 	"os"
// 	"crypto/ed25519"

// 	"riddhi/states/protocol/handler"
// 	"riddhi/states/protocol/packetlayer"
// 	"riddhi/states/protocol/statemachine"
// 	"riddhi/states/protocol/transportlayer"
// 	"riddhi/states/protocol/session"
// )

// func main() {

// 	// Connect to provider
// 	conn, err :=
// 		transportlayer.ConnectToProvider(
// 			"localhost:8080",
// 		)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer conn.Close()

// 	//create a output file to write the rebuilt file
// 	outputFile, err :=
// 		os.Create("rebuilt.pdf")

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer outputFile.Close()

// 	// Create a new state machine for the client
//     clientMachine :=
// 	statemachine.NewMachine(
// 		statemachine.Requester,
// 	)

// 	// Create a new session for the client
// 	clientSession := &session.Session{}
// 	fmt.Printf(
//     "client session addr: %p\n",
//      clientSession,
//     )

// 	//Generate signature
// 	clientPub, clientPriv, _ := ed25519.GenerateKey(nil)
//     //store signature
//     clientSession.RequesterPrivateKey = clientPriv
//     clientSession.RequesterPublicKey = clientPub

//     //request for fileiNfo
// 	transport := transportlayer.Transport{
//     Conn: conn,
//      }
// 	req :=
// 	packetlayer.FileMetadataRequestPayload{
// 		FileName: "sample.pdf",
// 	}

//     encoded :=
// 	packetlayer.EncodeFileMetadataRequest(
// 		req,
// 	)

//     packet :=
// 	packetlayer.Packet{
// 		Type: packetlayer.FileMetadataRequest,
// 		Length: uint32(len(encoded)),
// 		Payload: encoded,
// 	}

// 	err = transport.SendPacket(&packet)
//    if err != nil {
// 	panic(err)
//     }

//     responsePacket, err := transport.ReceivePacket()

//    if err != nil {
// 	panic(err)
//     }

//     metadata :=
// 	packetlayer.DecodeFileMetadataResponse(
// 		responsePacket.Payload,
// 	)

//     fmt.Printf(
// 	"received fileID: %x\n",
// 	metadata.FileID,
//      )

//     fmt.Println(
// 	"file size:",
// 	metadata.FileSize,
//      )

//     fmt.Println(
// 	"chunk size:",
//    	metadata.ChunkSize,
//      )

// 	 clientSession.FileID = metadata.FileID


// 	//nonce is for encryption(not signature)
// 	var nonce [16]byte
// 	copy(
// 		nonce[:],
// 		[]byte("abcdefghijklmnop"),
// 	)

// 	totalChunks := int(
// 		metadata.FileSize /
// 		uint64(metadata.ChunkSize),
// 	)

// for i := 0; i < totalChunks; i++ {

// 	req :=
// 		packetlayer.ChunkRequestPayload{
// 				ChunkIndex: uint32(i),
// 				Nonce:      nonce,
// 		}

// 	encodedReq :=
// 			packetlayer.EncodeChunkRequest(
// 				req,
// 		)

// 	requestPacket :=
// 			packetlayer.Packet{
// 				Type:
// 					packetlayer.ChunkRequest,

// 				Length:
// 					uint32(len(encodedReq)),

// 				Payload:
// 					encodedReq,
// 		}

// 	clientMachine.Transition(
// 			statemachine.WaitingForChunkResponse,
// 	    )

// 	fmt.Println(
// 			"Client requesting chunk",
// 			i,
// 		)

//     transport := transportlayer.Transport{
// 	Conn: conn,
//         }

//     err = transport.SendPacket(&requestPacket)
//     responsePacket, err :=
// 	transport.ReceivePacket()
// 	if err != nil {
// 		panic(err)
// 	}
// 	if responsePacket == nil {
// 		panic("nil response packet")
// 	}

//    //handle chunk response
// 	handler.HandleChunkResponse(
// 		responsePacket,
// 		clientMachine,
// 		clientSession,
// 	    )

//  for{

// 	ticketPacket :=
// 	handler.CreateLotteryTicket(
// 		clientSession,
// 	)
//   // SEND lottery ticket to provider
//     _,err =
// 	conn.Write(
// 		packetlayer.EncodePacket(
// 			byte(ticketPacket.Type),
// 			ticketPacket.Payload,
// 		),
// 	)

//     if err != nil {
// 	panic(err)
//     }

//   // RECEIVE key reveal from provider
//     revealPacket, err :=
// 	packetlayer.DecodePacket(
// 		conn,
// 	)

//     if err != nil {
// 	panic(err)
//     }

// 	//not winner, retry
// 	if revealPacket.Type != packetlayer.KeyReveal {
// 		fmt.Println(
// 			"ticket lost retrying",
// 		)

// 		// generate NEW ticket
// 		ticketPacket =
// 			handler.CreateLotteryTicket(
// 			clientSession,
// 		)

// 		continue
// 	}

// 	// winner, decrypt + verify
// 	if revealPacket.Type == packetlayer.KeyReveal {
//    // decrypt + verify
//     plaintext :=
// 	handler.HandleKeyReveal(
// 		revealPacket,
// 		clientMachine,
// 		clientSession,
// 	)

// 	if plaintext == nil {
// 		fmt.Println(
// 			"verification failed",
// 		)
// 		continue
// 	}

//     _, err =
// 	outputFile.Write(
// 		plaintext,
// 	)

//     if err != nil {
// 	panic(err)
//     }

// 	  break

// 	}

//   }


// }

// 	fmt.Printf( "client fileID: %x\n", clientSession.FileID, )

// 	fmt.Println(
// 		"file rebuilt",
// 	)

// }