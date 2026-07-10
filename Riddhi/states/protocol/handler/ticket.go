package handler

import (
	//"crypto/sha256"
	"fmt"

	"riddhi/states/protocol/packetlayer"
	"riddhi/states/protocol/session"
	"riddhi/states/protocol/statemachine"
	"riddhi/states/protocol/crypto"
)

//HandleLotteryTicket processes a received ticket
func HandleLotteryTicket(
	packet *packetlayer.Packet,
	machine *statemachine.Machine,
	session *session.Session,
) *packetlayer.Packet {


	//Validate role and state
	if machine.Role != statemachine.Provider {
	return nil
    }
	if machine.State !=
		statemachine.WaitingForLotteryTicket {

		fmt.Println("invalid state")
		return nil
	}

	//Decode ticket
	ticket :=
		packetlayer.DecodeLotteryTicket(
			packet.Payload,
		)

	//for debugging
	fmt.Println(
		"received lottery ticket for chunk:",
		session.CurrentChunkIndex,
	)
	//Store the requester's public key in the session
	session.RequesterPublicKey = ticket.ClientPublicKey

	//Verify signature
	valid :=
		 crypto.VerifyLotteryTicketSignature(
			session.RequesterPublicKey,
			ticket.ChunkIndex,
			ticket.TicketID,
			ticket.Signature,
		)

	if !valid {
		fmt.Println("invalid signature")
		return nil
	}


	// 0 = unpaid
	// 1 = paid

	isPaid := ticket.TicketID == 1

	if isPaid {
		fmt.Println(
			"Received PAID card for chunk",
			ticket.ChunkIndex,
		)

		// Later:
		// Verify smart contract payment
	} else {
		fmt.Println(
			"Received FREE card for chunk",
			ticket.ChunkIndex,
		)
	}


	//Build KeyReveal payload
	reveal :=
		packetlayer.KeyRevealPayload{
			ChunkIndex:
				ticket.ChunkIndex,

			Key: session.StoredKey,
		}

	//Encode reveal
	encoded :=
		packetlayer.EncodeKeyReveal(
			reveal,
		)

	// Wrap into packet
	responsePacket :=
		packetlayer.Packet{
			Type:
				packetlayer.KeyReveal,

			Length:
				uint32(len(encoded)),

			Payload:
				encoded,
		}

	// Transition state-reset for next chunk
	machine.Transition(
		statemachine.Idle,
	)

	return &responsePacket
}