package handler
import(
	"fmt"
	//"math/rand"
	"crypto/ed25519"
	"riddhi/states/protocol/packetlayer"
	"riddhi/states/protocol/statemachine"
	"riddhi/states/protocol/session"
	"riddhi/states/protocol/crypto"
)

func HandleChunkResponse(
	packet *packetlayer.Packet,
	machine *statemachine.Machine,
	session *session.Session,
) {

	//Validate role and state
	if machine.Role != statemachine.Requester{
	return 
      }
	if machine.State !=
		statemachine.WaitingForChunkResponse {

		fmt.Println("unexpected response")
		return 
	}


	// Decode response payload
	resp := packetlayer.DecodeChunkResponse(
		packet.Payload,
	)

	session.ProviderPublicKey =
	make([]byte, ed25519.PublicKeySize)

    copy(
	session.ProviderPublicKey,
	resp.ProviderPublicKey[:],
  )

    // Store the received data in the session
	session.FileID = resp.FileID
	session.CurrentChunkIndex = resp.ChunkIndex
	session.StoredCiphertext = resp.Ciphertext
    session.StoredNonce = resp.Nonce
    session.StoredCommitment = resp.Commitment
	session.Proof = resp.Proof
    session.MerkleRoot = resp.MerkleRoot

	//verify signature
	isValid := crypto.VerifyChunkResponseSignature(
		session.ProviderPublicKey ,
		resp.FileID,
		resp.ChunkIndex,
		resp.Commitment,
		resp.MerkleRoot,
		resp.ProviderSignature,
	)
	if !isValid {
		fmt.Println("invalid signature")
		return
	}

	// Store ciphertext
    fmt.Println(
	"received encrypted chunk:",
	len(resp.Ciphertext),
	"bytes",
    )
	//for debugging
	fmt.Println(
		"received nonce",
		len(resp.Nonce),
		"bytes",
	)
	
	// Store commitment hash
	fmt.Println(
	"received commitment",
)	
	
	//Update state
	machine.Transition(
		statemachine.WaitingForKeyReveal,
	)

	
}


// Ticket Generation and Sending to Provider
func CreateLotteryTicket(
	session *session.Session,
) *packetlayer.Packet{

//TicketID:= uint32(rand.Uint32())
//TicketID: uint32(time.Now().UnixNano()) //other way
//TicketID: uint32(i + 1000) // other way

paid := session.PaidChunks[session.CurrentChunkIndex]

var TicketID uint32

if paid {
	TicketID = 1
} else {
	TicketID = 0
}

// Generate signature
sig := crypto.SignLotteryTicket(
	session.RequesterPrivateKey,
	session.CurrentChunkIndex,
    TicketID,
)

ticket :=
	packetlayer.LotteryTicketPayload{
		ChunkIndex:session.CurrentChunkIndex,
		TicketID:TicketID,
		Signature:  sig,
		ClientPublicKey: session.RequesterPublicKey,
	}
message :=
	fmt.Sprintf(
		"%d:%d",
		ticket.ChunkIndex,
		ticket.TicketID,
	)
	//for debugging
	fmt.Println(
	"ticket message:",
	message,
)
	encoded :=
	packetlayer.EncodeLotteryTicket(
		ticket,
	)
	responsePacket :=
	packetlayer.Packet{
		Type:
			packetlayer.LotteryTicket,

		Length:
			uint32(len(encoded)),

		Payload:
			encoded,
	}

	return &responsePacket 
}