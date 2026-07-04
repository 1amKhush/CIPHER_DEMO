package packetlayer

import (
	"encoding/binary"
	"io"
	"net"
	"crypto/ed25519"
	"riddhi/states/protocol/crypto"
)

//encoder for filemetadata request
func EncodeFileMetadataRequest(
	req FileMetadataRequestPayload,
) []byte {

    nameBytes := []byte(req.FileName)

	payload :=
		make([]byte, 4+len(nameBytes))

	// filename length
	binary.BigEndian.PutUint32(
		payload[0:4],
		uint32(len(nameBytes)),
	)

	// filename bytes
	copy(
		payload[4:],
		nameBytes,
	)

	return payload
}

//decoder for filemetadata request
func DecodeFileMetadataRequest(
	data []byte,
) FileMetadataRequestPayload {

	var req FileMetadataRequestPayload

	nameLen :=
		binary.BigEndian.Uint32(
			data[0:4],
		)

	req.FileName =
		string(
			data[4:4+nameLen],
		)

	return req
}
//encoder for filemetadata response
func EncodeFileMetadataResponse(
	resp FileMetadataResponsePayload,
) []byte {

	payload := make([]byte, 44)

	offset := 0

	copy(
		payload[offset:offset+32],
		resp.FileID[:],
	)

	offset += 32

	binary.BigEndian.PutUint64(
		payload[offset:offset+8],
		resp.FileSize,
	)

	offset += 8

	binary.BigEndian.PutUint32(
		payload[offset:offset+4],
		resp.ChunkSize,
	)

	offset += 4

	return payload
}

//decoder for filemetadata response
func DecodeFileMetadataResponse(
	data []byte,
) FileMetadataResponsePayload {

	var resp FileMetadataResponsePayload

	offset := 0

	copy(
		resp.FileID[:],
		data[offset:offset+32],
	)

	offset += 32

	resp.FileSize =
		binary.BigEndian.Uint64(
			data[offset:offset+8],
		)

	offset += 8

	resp.ChunkSize =
		binary.BigEndian.Uint32(
			data[offset:offset+4],
		)

	offset += 4

	return resp

}

//encoder-for packet
func EncodePacket(t byte, payload []byte) []byte {
	packet := make([]byte, 1+4+len(payload))

	packet[0] = t

	binary.BigEndian.PutUint32(
		packet[1:5],
		uint32(len(payload)),
	)

	copy(packet[5:], payload)

	return packet
}

//decoder-for packet
func DecodePacket(conn net.Conn) (*Packet, error) {

	typeBuf := make([]byte, 1)

	_, err := io.ReadFull(conn, typeBuf)
	if err != nil {
		return nil, err
	}

	lenBuf := make([]byte, 4)

	_, err = io.ReadFull(conn, lenBuf)
	if err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint32(lenBuf)

	payload := make([]byte, size)

	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return nil, err
	}

	return &Packet{
		Type:    PacketType(typeBuf[0]),//don't use typebuf[0],both are byte but Go treat PacketType and typebuf differently
		Length:  size,
		Payload: payload,
	}, nil
}

//for chunk
func EncodeChunkRequest(
	req ChunkRequestPayload,
) []byte {

	payload := make([]byte, 20)

	// ChunkIndex → first 4 bytes
	binary.BigEndian.PutUint32(
		payload[0:4],
		req.ChunkIndex,
	)

	// Nonce → next 16 bytes
	copy(payload[4:20], req.Nonce[:])

	return payload
}

func DecodeChunkRequest(
	data []byte,
) ChunkRequestPayload {

	var req ChunkRequestPayload

	req.ChunkIndex = binary.BigEndian.Uint32(
		data[0:4],
	)
	
	copy(req.Nonce[:], data[4:20]) //Nonce[16]byte is fixed size array

	return req
}


//encoder of chunkresponse
func EncodeChunkResponse(
	resp ChunkResponsePayload,
) []byte {

	cipherLen := len(resp.Ciphertext)
	proofCount := len(resp.Proof.Steps)
	total := 32 + // FileID 
	 4 + // ChunkIndex 
	  24 + // Nonce 
	  4 + // CipherLen 
	 cipherLen + 
	 32 + // Commitment 
	  32 + // MerkleRoot 
	  4 + // ProofCount
	  (proofCount * 33) + // ProofSteps
	  64 +// ProviderSignature
      ed25519.PublicKeySize // ProviderPublicKey

	payload := make(
		[]byte,
		total,
	)

	offset := 0 
	copy(payload[offset:offset+32], resp.FileID[:]) 
	offset += 32

	binary.BigEndian.PutUint32( //for integer you should use(for bytes)
	payload[offset:offset+4],
	resp.ChunkIndex,
)

    offset += 4
	
	// Don't forget to copy nonce
	copy(
		payload[offset:offset+24],
		resp.Nonce[:],
	)

	offset += 24

	// ciphertext length
	binary.BigEndian.PutUint32(
		payload[offset:offset+4],
		uint32(cipherLen),
	)

	offset+=4
	
	// Ciphertext bytes
	copy(
		payload[offset:offset+cipherLen],
		resp.Ciphertext,
	)

	offset += cipherLen

	// Commitment bytes
	copy(
		payload[offset:],
		resp.Commitment[:],
	)

	offset += 32 
	copy( 
		payload[offset:offset+32],
	    resp.MerkleRoot[:], 
		) 
		offset += 32 
		
		binary.BigEndian.PutUint32( 
			payload[offset:offset+4], 
			uint32(proofCount), ) 
			
			offset += 4 
			
			for _, step:= range resp.Proof.Steps {
				 copy( 
					payload[offset:offset+32], 
					step.Hash[:], 
					) 
					offset += 32 
					if step.Left {
						payload[offset] = 1
					} else {
						payload[offset] = 0
					}
					offset += 1
				}
		copy(
			payload[offset:offset+64],
			resp.ProviderSignature[:],
		)

		offset += 64

	// ProviderPublicKey bytes
	copy(
		payload[offset:offset+ed25519.PublicKeySize],
		resp.ProviderPublicKey,
	)
	
	offset += ed25519.PublicKeySize

	return payload
}


//decoder for chunk response
func DecodeChunkResponse(
	data []byte,
) ChunkResponsePayload {

	var resp ChunkResponsePayload

	offset := 0

	copy( 
		resp.FileID[:], 
		data[offset:offset+32], 
	)

	offset += 32

	// STORE the chunk index
	resp.ChunkIndex =
		binary.BigEndian.Uint32(
			data[offset:offset+4],
		)

		offset+=4


	//don't forget to copy
	copy(
		resp.Nonce[:],
		data[offset:offset+24],
	)

	offset += 24

	// Read ciphertext length
	cipherLen := binary.BigEndian.Uint32(
		data[offset:offset+4],
	)

	offset += 4

	// Extract ciphertext
	resp.Ciphertext = make(
		[]byte,
		cipherLen,
	)

	copy(
		resp.Ciphertext,
		data[offset:offset+int(cipherLen)],
	)

	offset += int(cipherLen)

	// Extract Commitment
	copy( resp.Commitment[:], data[offset:offset+32], ) 
	offset += 32

	copy( resp.MerkleRoot[:], data[offset:offset+32], ) 
	offset += 32

	proofCount := binary.BigEndian.Uint32( data[offset:offset+4], ) 
	offset += 4

	resp.Proof.Steps = 
	   make([]crypto.ProofStep, proofCount)

	   for i := 0; i < int(proofCount); i++ { 
		copy(
			 resp.Proof.Steps[i].Hash[:], 
			data[offset:offset+32], ) 
			offset += 32
			resp.Proof.Steps[i].Left = data[offset] == 1
			offset += 1
		}

		copy(
			resp.ProviderSignature[:],
			data[offset:offset+64],
		)

		offset += 64

	// Extract ProviderPublicKey
	resp.ProviderPublicKey =
	make([]byte, ed25519.PublicKeySize)

  copy(
	resp.ProviderPublicKey,
	data[offset:offset+ed25519.PublicKeySize],
   )

  offset += ed25519.PublicKeySize
	
	return resp
}

//encoder for lotteryticket
func EncodeLotteryTicket(
	ticket LotteryTicketPayload,
) []byte {

	payload := make([]byte, 72 + ed25519.PublicKeySize)

	binary.BigEndian.PutUint32(
		payload[0:4],
		ticket.ChunkIndex,
	)

	binary.BigEndian.PutUint32(
		payload[4:8],
		ticket.TicketID,
	)

	copy(
		payload[8:72],
		ticket.Signature[:],
	)

	copy(
		payload[72:72+ed25519.PublicKeySize],
		ticket.ClientPublicKey,
	)

	return payload
}

//decoder for lottery ticket
func DecodeLotteryTicket(
	data []byte,
) LotteryTicketPayload {

	var ticket LotteryTicketPayload

	ticket.ChunkIndex =
		binary.BigEndian.Uint32(
			data[0:4],
		)

	ticket.TicketID =
		binary.BigEndian.Uint32(
			data[4:8],
		)

	copy(
		ticket.Signature[:],
		data[8:72],
	)
	ticket.ClientPublicKey =
	make([]byte, ed25519.PublicKeySize)

copy(
	ticket.ClientPublicKey,
	data[72:72+ed25519.PublicKeySize],
)

	return ticket
}

//encoder for key
func EncodeKeyReveal(
	key KeyRevealPayload,
) []byte {

	payload := make([]byte, 36)

	binary.BigEndian.PutUint32(
		payload[0:4],
		key.ChunkIndex,
	)

	copy(
		payload[4:36],////cannot convert whole byte array to uint32 You should COPY raw bytes.
		key.Key[:],
	)

	return payload
}
//decoder for key
func DecodeKeyReveal(
	data []byte,
) KeyRevealPayload {

	var key KeyRevealPayload

	key.ChunkIndex =
		binary.BigEndian.Uint32(
			data[0:4],
		)

	copy(
		key.Key[:],
		data[4:36],
	)

	return key
}

//encoder for ticket reject
func EncodeTicketReject(
	reject TicketRejectPayload,
) []byte {

	payload := make([]byte, 1)

	payload[0] = reject.Reason

	return payload
}

//decoder for ticket reject
func DecodeTicketReject(
	data []byte,
) TicketRejectPayload {

	return TicketRejectPayload{
		Reason: data[0],
	}
}