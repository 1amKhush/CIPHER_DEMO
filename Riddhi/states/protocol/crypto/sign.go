package crypto

import(
	 "crypto/ed25519"
	 "fmt"
)

//provider signing function
func SignChunkResponse(
privateKey ed25519.PrivateKey,
fileID[32]byte,
chunkIndex uint32,
commitment[32]byte,
merkleRoot[32]byte,
)[64]byte{
	message := fmt.Sprintf(
		 "%x:%d:%x:%x", 
		 fileID, 
		 chunkIndex,
		 commitment, 
		 merkleRoot, 
		 ) 
		 signature := ed25519.Sign( 
			privateKey, 
			[]byte(message),
			 ) 
			 var sig [64]byte
			  copy(sig[:], signature)

			 return sig

			}
			
//verification of provider signature
func VerifyChunkResponseSignature(
publicKey ed25519.PublicKey,
fileID [32]byte,
chunkIndex uint32,
commitment [32]byte,
merkleRoot [32]byte,
signature [64]byte,
) bool {


message := fmt.Sprintf(
	"%x:%d:%x:%x",
	fileID,
	chunkIndex,
	commitment,
	merkleRoot,
)

return ed25519.Verify(
	publicKey,
	[]byte(message),
	signature[:],
)

}
	
// requester signing function
func SignLotteryTicket(
privateKey ed25519.PrivateKey,
chunkIndex uint32,
ticketID uint32,
) [64]byte {

message := fmt.Sprintf(
	"%d:%d",
	chunkIndex,
	ticketID,
)

signature := ed25519.Sign(
	privateKey,
	[]byte(message),
)

var sig [64]byte

copy(
	sig[:],
	signature,
)

return sig

}

// verification of requester signature
func VerifyLotteryTicketSignature(
publicKey ed25519.PublicKey,
chunkIndex uint32,
ticketID uint32,
signature [64]byte,
) bool {

message := fmt.Sprintf(
	"%d:%d",
	chunkIndex,
	ticketID,
)

return ed25519.Verify(
	publicKey,
	[]byte(message),
	signature[:],
)

}
