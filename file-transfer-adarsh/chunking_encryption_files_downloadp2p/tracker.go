package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

// central phonebook that remembers which peer has what file

// TrackerMessage is the protocol used between peers and the tracker
type TrackerMessage struct {
	Type     string `json:"type"`     // "ANNOUNCE", "QUERY", "RESPONSE"
	FileName string `json:"filename"` 
	PeerAddr string `json:"peerAddr"` 
}

// Global registry: maps filename to a list of peer addresses
var fileRegistry = make(map[string][]string)
var registryMutex = sync.Mutex{}

func RunTrackerServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Tracker listen error:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Tracker Server running on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Tracker accept error:", err)
			continue
		}
		
		go handleTrackerConnection(conn)
	}
}

func handleTrackerConnection(conn net.Conn) {
	defer conn.Close()

	payload, err := receivePayload(conn)
	if err != nil {
		fmt.Println("Error receiving tracker payload:", err)
		return
	}

	var msg TrackerMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		fmt.Println("Error parsing tracker message:", err)
		return
	}

	registryMutex.Lock()
	defer registryMutex.Unlock()

	switch msg.Type {
	case "ANNOUNCE":
		peers := fileRegistry[msg.FileName]
		// Add if not exists
		exists := false
		for _, p := range peers {
			if p == msg.PeerAddr {
				exists = true
				break
			}
		}
		if !exists {
			fileRegistry[msg.FileName] = append(fileRegistry[msg.FileName], msg.PeerAddr)
		}
		fmt.Printf("Tracker: Registered %s for file %s\n", msg.PeerAddr, msg.FileName)
		
		// Send ACK
		resp := TrackerMessage{Type: "RESPONSE", PeerAddr: "OK"}
		respBytes, _ := json.Marshal(resp)
		sendPayload(conn, respBytes)

	case "QUERY":
		fmt.Printf("Tracker: Received query for file %s\n", msg.FileName)
		peers := fileRegistry[msg.FileName]
		peerAddr := ""
		if len(peers) > 0 {
			peerAddr = peers[0] // Simple: return the first available peer
		}
		
		resp := TrackerMessage{Type: "RESPONSE", FileName: msg.FileName, PeerAddr: peerAddr}
		respBytes, _ := json.Marshal(resp)
		sendPayload(conn, respBytes)
	}
}
