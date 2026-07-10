package handler

import (
	
	"fmt"

	"riddhi/states/protocol/crypto"
	"riddhi/states/protocol/packetlayer"
	"riddhi/states/protocol/statemachine"
	"riddhi/states/protocol/session"
)

//HandleKeyReveal processes a received key reveal
func HandleKeyReveal(
	packet *packetlayer.Packet,
	machine *statemachine.Machine,
	session *session.Session,
)[]byte{

	//Validate role and state
	if machine.Role != statemachine.Requester {
	return nil
    }
	if machine.State !=
		statemachine.WaitingForKeyReveal {

		fmt.Println(
			"invalid state",
		)

		return nil
	}

	//Decode reveal
	reveal :=
		packetlayer.DecodeKeyReveal(
			packet.Payload,
		)

	//Verify commitment
	valid :=
		crypto.VerifyCommitment(
			reveal.Key,
			session.StoredCiphertext,
		    session.StoredCommitment,
		)

	if !valid {

		fmt.Println(
			"commitment verification failed",
		)

		return nil
	}

	fmt.Println(
		"commitment verified",
	)

	//for debugging
	fmt.Printf(
	"client nonce: %x\n",
	session.StoredNonce,
    )
   fmt.Printf(
	"client key: %x\n",
	reveal.Key,
    ) 

	//Decrypt chunk
	plaintext, err :=
		crypto.Decrypt(
			reveal.Key,
		session.StoredNonce,
		session.StoredCiphertext,
		)

	if err != nil {

		fmt.Println(
			"decryption failed",
		)

		return nil
	}

	root:= session.MerkleRoot
	index:= int(session.CurrentChunkIndex)
	chunk:= plaintext
	fileID:= session.FileID
	proof:= session.Proof

	leaf := crypto.GenerateLeafHash( 
		fileID, index, len(chunk),chunk,
	)
	//for debugging
	fmt.Printf(
	"client leaf: %x\n",
	leaf,
  )

	//Verify plaintext
	if crypto.VerifyProof(
		 root, leaf, proof, 
		){
	fmt.Println(
		"chunk verified successfully",
	)
    	} else {
			fmt.Println(
				"chunk verification failed",
			)
			return nil
		}

    //for debugging
    fmt.Println(
	"decrypted bytes:",
	len(plaintext),
   )

	// Reset protocol
	machine.Transition(
		statemachine.Idle,
	)

	return plaintext
}