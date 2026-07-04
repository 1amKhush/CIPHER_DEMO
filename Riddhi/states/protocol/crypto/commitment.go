package crypto
import(
	//"crypto/rand"
)
// type CommitmentData struct {
// 	Secret     [32]byte //we can add secret which belongs to provider and is used to generate commitment
// 	Commitment [32]byte
// }

func GenerateCommitment(
	key [32]byte,
	ciphertext []byte,
)([32]byte,error){

	// var result CommitmentData

	// // random secret
	// _, err := rand.Read(
	// 	result.Secret[:],
	// )

	// if err != nil {
	// 	return result, err
	// }

	combined := make(
		[]byte,
		0,
		32+len(ciphertext),
	)

	combined = append(
		combined,
		key[:]...,
	)

	combined = append(
		combined,
		ciphertext...,
	)

	result:=
		GenerateHash(combined)

	return result,nil
}

func VerifyCommitment(
	// secret [32]byte,
	// data []byte,
	// commitment [32]byte,
	key [32]byte,
	ciphertext []byte,
	commitment [32]byte,
) bool {

	combined := make(
		[]byte,
		0,
		32+len(ciphertext),
	)

	combined = append(
		combined,
		key[:]...,
	)

	combined = append(
		combined,
		ciphertext...,
	)

	computed :=
		GenerateHash(combined)

	return computed == commitment
}