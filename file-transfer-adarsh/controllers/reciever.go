package controllers

import (
	// "bufio"
	"fmt"
	"io"
	"net"
	"os"
)




func Reciever() {
	fmt.Println("You are on port: 8090 to recieve file ")

	server, err := net.Listen("tcp", ":8090")
	if err != nil {
		fmt.Println("Error listening on port 8090:", err)
		return
	}
	defer server.Close()

	fmt.Println("Server is listening on port 8090...")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())

	// Read everything from connection
	data, err := io.ReadAll(conn)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	method, fname, payload, err := AnalyseRequest(data)
	if err != nil {
		fmt.Println("Invalid protocol:", err)
		conn.Write([]byte("STATUS 400\n\nMalformed request"))
		return
	}

	if method == "POST" {
		// Ensure the output directory exists
		err = os.MkdirAll("output", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		
		// We no longer save to the output directory to prevent disk space leaks.
		// The Web UI handles the actual file download.
		fmt.Printf("TCP Receiver verified file %s successfully! Bytes received: %d\n", fname, len(payload))
		conn.Write([]byte("STATUS 200\n\nOK"))
	} else if method == "GET" {
		// GET is not supported by the dummy receiver anymore, 
		// because files are stored in staging/ and managed by the Web API.
		conn.Write([]byte("STATUS 404\n\nFile Not Found"))
	} else {
		conn.Write([]byte("STATUS 400\n\nUnknown Method"))
	}
}
