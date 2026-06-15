package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// ======================
// Protocol Types
// ======================

type Request struct {
	Method string
	Path   string
	Size   int
	Body   []byte
}

type Response struct {
	Status int
	Size   int
	Body   []byte
}

// ======================
// Protocol Parser
// ======================

func ReadRequest(reader *bufio.Reader) (*Request, error) {

	req := &Request{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "METHOD":
			req.Method = value

		case "PATH":
			req.Path = value

		case "SIZE":
			size, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			req.Size = size
		}
	}

	// Read body if present
	if req.Size > 0 {
		req.Body = make([]byte, req.Size)

		_, err := io.ReadFull(reader, req.Body)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

// ======================
// Protocol Serializer
// ======================

func WriteResponse(conn net.Conn, resp Response) error {
	resp.Size = len(resp.Body)
	header := fmt.Sprintf( "STATUS:%d\nSIZE:%d\n\n", resp.Status, resp.Size)

	_, err := conn.Write([]byte(header))
	if err != nil {
		return err
	}

	if resp.Size > 0 {
		_, err = conn.Write(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

// ======================
// GET Handler
// ======================

func handleGet(conn net.Conn, req *Request) {
	path := "./Server/DataBase" + req.Path
	data, err := os.ReadFile(path)

	if err != nil {
		WriteResponse(conn, Response{
			Status: 404,
			Body:   []byte("File not found"),
		})
		return
	}

	WriteResponse(conn, Response{
		Status: 200,
		Body:   data,
	})
}

// ======================
// POST Handler
// ======================

func handlePost(conn net.Conn, req *Request) {
	path := "./Server/DataBase" + req.Path

	err := os.WriteFile(path,req.Body,0644)

	if err != nil {
		WriteResponse(conn, Response{
			Status: 500,
			Body:   []byte("Failed to save file"),
		})
		return
	}

	WriteResponse(conn, Response{
		Status: 200,
		Body:   []byte("File stored successfully"),
	})
}

// ======================
// List Handler
// ======================
func handleList(conn net.Conn) {
	files, err := os.ReadDir("./Server/DataBase")

	if err != nil {
		WriteResponse(conn, Response{
			Status: 500,
			Body:   []byte("Error"),
		})
		return
	}

	var names string

	for _, file := range files {

		if !file.IsDir() {
			names += file.Name() + "\n"
		}
	}

	WriteResponse(conn, Response{
		Status: 200,
		Body:   []byte(names),
	})
}

// ======================
// Connection Handler
// ======================

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	req, err := ReadRequest(reader)

	if err != nil {
		fmt.Println("Request error:", err)
		WriteResponse(conn, Response{
			Status: 400,
			Body:   []byte("Bad request"),
		})
		return
	}

	fmt.Println("================================")
	fmt.Println("Method:", req.Method)
	fmt.Println("Path:", req.Path)
	fmt.Println("Size:", req.Size)
	fmt.Println("================================")

	switch req.Method {

	case "GET":
		handleGet(conn, req)

	case "POST":
		handlePost(conn, req)

	case "LIST":
		handleList(conn)

	default:

		WriteResponse(conn, Response{
			Status: 400,
			Body:   []byte("Unknown method"),
		})
	}
}

// ======================
// Server Entry Point
// ======================

func main() {

	// Create a DataBase folder if not present
	os.MkdirAll("./Server/DataBase", os.ModePerm)

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}

	fmt.Println("================================")
	fmt.Println("Custom TCP File Server Running")
	fmt.Println("Listening on :8080")
	fmt.Println("================================")

	// This is my infinite loop for listening the requests from the cliet

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		fmt.Println("Client Connected:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}