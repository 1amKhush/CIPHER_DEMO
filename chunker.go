package main

const ChunkSize = 32768

func SplitIntoChunks(data []byte) [][]byte {
	var chunks [][]byte

	for start := 0; start < len(data); start += ChunkSize {
		end := start + ChunkSize

		if end > len(data) {
			end = len(data)
		}
        //bytes slicing of 32 kilo bytes
		chunks = append(chunks, data[start:end])
	}

	return chunks
}