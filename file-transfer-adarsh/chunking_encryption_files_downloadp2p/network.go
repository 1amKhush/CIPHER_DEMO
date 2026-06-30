package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

// does all the tcp network talking for tracker and file transfer

// Transport defines the interface for our network protocol
// allowing us to swap out TCP for QUIC in the future.
type Transport interface {
	ListenForRequests(port string, handler func(filename string) (Metadata, KeyNonceMap, [][]byte, error)) error
	RequestFile(addr string, filename string) (Metadata, KeyNonceMap, [][]byte, error)
}

type TCPTransport struct{}

func (t *TCPTransport) ListenForRequests(port string, handler func(filename string) (Metadata, KeyNonceMap, [][]byte, error)) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("Listening for P2P requests on TCP port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		
		go func(c net.Conn) {
			defer c.Close()
			
			// Receive request string
			reqBytes, err := receivePayload(c)
			if err != nil {
				fmt.Println("Receive error:", err)
				return
			}
			filename := string(reqBytes)
			fmt.Println("Requested file:", filename)
			
			metadata, keys, chunks, err := handler(filename)
			if err != nil {
				fmt.Println("Handler error:", err)
				// Send empty metadata to signal error
				sendPayload(c, []byte("{}"))
				return
			}
			
			// 1. Send Metadata
			metaBytes, _ := json.Marshal(metadata)
			sendPayload(c, metaBytes)

			// 2. Send Keys and Nonces
			keysBytes, _ := json.Marshal(keys)
			sendPayload(c, keysBytes)

			// 3. Send Chunks
			for _, chunk := range chunks {
				sendPayload(c, chunk)
			}
		}(conn)
	}
}

func (t *TCPTransport) RequestFile(addr string, filename string) (Metadata, KeyNonceMap, [][]byte, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return Metadata{}, nil, nil, err
	}
	defer conn.Close()

	// Send file request
	if err := sendPayload(conn, []byte(filename)); err != nil {
		return Metadata{}, nil, nil, err
	}

	var metadata Metadata
	var keys KeyNonceMap
	var chunks [][]byte

	// 1. Receive Metadata
	metaBytes, err := receivePayload(conn)
	if err != nil {
		return metadata, nil, nil, err
	}
	if string(metaBytes) == "{}" {
		return metadata, nil, nil, fmt.Errorf("file not found on peer")
	}
	json.Unmarshal(metaBytes, &metadata)

	// 2. Receive Keys and Nonces
	keysBytes, err := receivePayload(conn)
	if err != nil {
		return metadata, nil, nil, err
	}
	json.Unmarshal(keysBytes, &keys)

	// 3. Receive Chunks
	for i := 0; i < metadata.ChunkCount; i++ {
		chunk, err := receivePayload(conn)
		if err != nil {
			return metadata, keys, nil, err
		}
		chunks = append(chunks, chunk)
	}

	return metadata, keys, chunks, nil
}

// Signaling functions
func AnnounceToTracker(trackerAddr, filename, myListenAddr string) error {
	conn, err := net.Dial("tcp", trackerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := TrackerMessage{Type: "ANNOUNCE", FileName: filename, PeerAddr: myListenAddr}
	msgBytes, _ := json.Marshal(msg)
	
	if err := sendPayload(conn, msgBytes); err != nil {
		return err
	}

	respBytes, err := receivePayload(conn)
	if err != nil {
		return err
	}
	var resp TrackerMessage
	json.Unmarshal(respBytes, &resp)
	if resp.PeerAddr != "OK" {
		return fmt.Errorf("tracker announcement failed")
	}
	return nil
}

func QueryTracker(trackerAddr, filename string) (string, error) {
	conn, err := net.Dial("tcp", trackerAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	msg := TrackerMessage{Type: "QUERY", FileName: filename}
	msgBytes, _ := json.Marshal(msg)
	
	if err := sendPayload(conn, msgBytes); err != nil {
		return "", err
	}

	respBytes, err := receivePayload(conn)
	if err != nil {
		return "", err
	}
	var resp TrackerMessage
	json.Unmarshal(respBytes, &resp)
	
	if resp.PeerAddr == "" {
		return "", fmt.Errorf("no peers available for file %s", filename)
	}
	return resp.PeerAddr, nil
}

// helper to send length-prefixed data
func sendPayload(conn net.Conn, data []byte) error {
	length := uint32(len(data))
	lenBuf := []byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}
	if _, err := conn.Write(lenBuf); err != nil {
		return err
	}
	_, err := conn.Write(data)
	return err
}

func receivePayload(conn net.Conn) ([]byte, error) {
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return nil, err
	}
	length := uint32(lenBuf[0])<<24 | uint32(lenBuf[1])<<16 | uint32(lenBuf[2])<<8 | uint32(lenBuf[3])
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}
	return data, nil
}
