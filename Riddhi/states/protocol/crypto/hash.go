package crypto

import (
"encoding/binary"
"bytes"
"fmt"
 "os"
	
	"golang.org/x/crypto/sha3"
)

//general hash
func GenerateHash(data []byte)[32]byte {
  // Generate a unique hash for the chunk data
  hash :=sha3.NewLegacyKeccak256()
  hash.Write(data) 
  result:=hash.Sum(nil) 
  
  //keep it raw bytes instead of string
  //return hex.EncodeToString(result)
  
  var finalHash [32]byte

  copy(finalHash[:], result)

  return finalHash
}
//hash for file
func GenerateHashFile(
	path string,
) [32]byte {

	data, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return GenerateHash(data)
}

//leaf hash
func GenerateLeafHash(
	fileID [32]byte,
	index int,
	length int,
	chunk []byte,
) [32]byte {

	var buf bytes.Buffer

	// optional prefix for leaf nodes
	buf.WriteByte(0x00)

	// FileID
	buf.Write(fileID[:])

	// Chunk index
	err:=binary.Write(&buf, binary.BigEndian, uint32(index))
	if err != nil {
		panic(err)
	}
	// Chunk length 
	err = binary.Write(&buf, binary.BigEndian, uint32(len(chunk)), ) 
	if err != nil {
		 panic(err)
		 }

	// Raw chunk bytes
	buf.Write(chunk)

	// Final leaf hash
	return GenerateHash(buf.Bytes())
}

//parent hash
func GenerateParentHash(
	left [32]byte,
	right [32]byte,
) [32]byte {


	var buf bytes.Buffer

	// optional prefix for internal nodes
	buf.WriteByte(0x01)

	// Left child hash
	buf.Write(left[:])

	// Right child hash
	buf.Write(right[:])

	// Final parent hash
	return GenerateHash(buf.Bytes())
}

//generate root
func BuildMerkleRoot(
	leaves [][32]byte,
	) [32]byte {

	 if len(leaves) == 0 {
		 return [32]byte{} 
		 } 

		current := leaves

		for len(current) > 1 {

		// If the number of nodes is odd, duplicate the last node
		  if len(current)%2 != 0 { 
		    current = append(current, current[len(current)-1]) 
		}

		// Create a new slice to hold the parent hashes
		next := make([][32]byte, 0, len(current)/2)

	    for i := 0; i < len(current); i += 2 {

		   parent := GenerateParentHash(current[i], current[i+1])

		   //append means add the parent hash to the next level of the tree
		   next = append(next, parent) 
		   } 

		   // Move to the next level of the tree
		   current = next    
		 } 

		   return current[0]

 }

//Generate proof
type ProofStep struct { 
	Hash [32]byte 
	Left bool  //left=true, right=false
}

type MerkleProof struct { 
	Steps []ProofStep 
}

func GenerateProof( 
	leaves [][32]byte, targetIndex int,
	 ) (MerkleProof, error) { 
		if targetIndex < 0 || targetIndex >= len(leaves) { 
			return MerkleProof{}, fmt.Errorf("invalid target index") 
		} 
	var proof MerkleProof 

	current := leaves 
	index := targetIndex 
	
	for len(current) > 1 { 

		// Handle odd nodes 
		if len(current)%2 != 0 { 
			current = append(current, current[len(current)-1]) 
		} 

		var next [][32]byte
		
		for i := 0; i < len(current); i += 2 { 
			
			left := current[i] 
			right := current[i+1] 
			
			parent := GenerateParentHash(left, right) 
			next = append(next, parent) 
			
			// check if the current index is either i or i+1, which means the target is in this pair
			if i == index || i+1 == index { 

				if index == i { 
					// sibling on right 
					proof.Steps = append( 
						proof.Steps, 
						ProofStep{ 
							Hash: right, 
							Left: false,
						 }, 
					) 
			} else if(index == i+1) {
				 // sibling on left 
				 proof.Steps = append( 
					proof.Steps, 
					ProofStep{ 
						Hash: left,
					    Left: true, 
					    }, 
					) 
			} 
				// update index for the next level
				index = len(next) - 1 
				} 
		} 

		current = next 

	} 

			return proof, nil

 }

// Verify proof
func VerifyProof(
	root [32]byte, 
	leaf [32]byte,
    proof MerkleProof,
	) bool{
	
		current:=leaf

		for _,step:=range proof.Steps{

			if step.Left{
				current = GenerateParentHash(
					step.Hash,
					current,
				)
			}else{
				current =GenerateParentHash(
					current,
					step.Hash,
				)
			}
		}
		return current == root

}
	