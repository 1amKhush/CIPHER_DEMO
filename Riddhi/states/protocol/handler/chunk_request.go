package handler

import(
	"fmt"
	"riddhi/states/protocol/packetlayer"
	"riddhi/states/protocol/statemachine"
	"riddhi/states/protocol/crypto"
	"riddhi/states/protocol/session"
	"riddhi/states/storage"
)

//handlemetainfo request:-
func HandleFileMetadataRequest(
	packet *packetlayer.Packet,
) *packetlayer.Packet {

	req :=
		packetlayer.DecodeFileMetadataRequest(
			packet.Payload,
		)

	//for debugging
	fmt.Println(
		"metadata requested for:",
		req.FileName,
	)

	resp :=
		packetlayer.FileMetadataResponsePayload{
			FileID:
				storage.ChunkedFile.FileID,

			FileSize:
				uint64(
					storage.ChunkedFile.FileSize,
				),

			ChunkSize:
				uint32(
					storage.ChunkedFile.ChunkSize,
				),

		}

	encoded :=
		packetlayer.EncodeFileMetadataResponse(
			resp,
		)

	responsePacket :=
		packetlayer.Packet{
			Type:
				packetlayer.FileMetadataResponse,

			Length:
				uint32(len(encoded)),

			Payload:
				encoded,
		}

	return &responsePacket
}

// HandleChunkRequest handles the chunk request from the requester
func HandleChunkRequest(
	packet *packetlayer.Packet,
	machine *statemachine.Machine,
	session *session.Session,
)*packetlayer.Packet {

	// Validate role and state 
	if machine.Role != statemachine.Provider{
	return nil
    }
	if machine.State!= statemachine.Idle {
		fmt.Println("invalid state")
		return nil
	}

	// Decode payload
	req := packetlayer.DecodeChunkRequest(
		packet.Payload,
	)

	// Process logic
	fmt.Println("requested chunk:",
		req.ChunkIndex)


	if int(req.ChunkIndex) >= len(storage.ChunkedFile.Chunks) {
	fmt.Println("invalid chunk index")
	return nil
    }

    //locate chunk
    requestedChunk :=
	storage.ChunkedFile.Chunks[req.ChunkIndex]

	// actual chunk bytes
    data := requestedChunk.ChunkData

	// Encrypt chunk

    //lets make key first
    var key [32]byte
    copy(
	key[:],
	[]byte(
		"abcdefghijklmnopqrstuvwxyz123456",
	    ),
    )

	//for debugging
    fmt.Printf(
	 "provider key: %x\n",
	 key,
    )

    //use it 
    nonce, ciphertext, err := crypto.Encrypt(data, key)
    if err != nil {
	  return nil
    }

	//store everything
	session.CurrentChunkIndex = req.ChunkIndex
	session.StoredCiphertext = ciphertext	
    session.StoredNonce = nonce
    session.StoredKey = key
    session.OriginalChunk = data
	
	fmt.Printf( "provider leaf: %x\n", requestedChunk.LeafHash, )
	// generate proof here 
	proof, err := crypto.GenerateProof(storage.ChunkedFile.LeafHashes, int(req.ChunkIndex)) 
	if err != nil {
		 return nil 
		}

	// Generate commitment
	Commitment,err :=
		crypto.GenerateCommitment(key,ciphertext)
		if err!=nil{
			return nil
		}

		//store
		session.StoredCommitment = Commitment

		sign:=crypto.SignChunkResponse(
			session.ProviderPrivateKey,
			storage.ChunkedFile.FileID,
			session.CurrentChunkIndex,
			Commitment,
			storage.ChunkedFile.MerkleRoot,
		)
    
	// Build ChunkResponse
	resp :=
		packetlayer.ChunkResponsePayload{
			FileID: storage.ChunkedFile.FileID,
			ChunkIndex: session.CurrentChunkIndex ,
			Nonce: nonce,
			Ciphertext: ciphertext,
			Commitment: Commitment,
			Proof: proof,
			MerkleRoot: storage.ChunkedFile.MerkleRoot,
			ProviderSignature:sign,
		    ProviderPublicKey: session.ProviderPublicKey,
		}

	encoded :=
		packetlayer.EncodeChunkResponse(resp)

	// wrap into packet
	responsePacket :=
		packetlayer.Packet{
			Type:packetlayer.ChunkResponse,
			Length: uint32(len(encoded)),
			Payload: encoded,
		}

	fmt.Println(
		"generated response packet",
		responsePacket.Type,
	)	

	// Update state
	machine.Transition(
		statemachine.WaitingForLotteryTicket,
	)
	return &responsePacket
}



