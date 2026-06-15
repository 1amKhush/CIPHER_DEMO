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

func (s *Server) Handleconnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	trimmed := strings.TrimSpace(line)
	parts := strings.Fields(trimmed)
	if len(parts) < 2 {
		fmt.Fprintf(conn, "STATUS 400\r\nLength: 0\r\n\r\n")
		return
	}
	method := parts[0]
	path := parts[1]

	if method == "GET" {
		req := Request{
			Method: method,
			Path:   path,
			Body:   nil,
		}
		GetFile(conn, req)
	} else if method == "POST" {
		if len(parts) < 3 {
			fmt.Fprintf(conn, "STATUS 400\r\nLength: 0\r\n\r\n")
			return
		}
		contentLength, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Fprintf(conn, "STATUS 400\r\nLength: 0\r\n\r\n")
			return
		}
		body := make([]byte, contentLength)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			fmt.Fprintf(conn, "STATUS 500\r\nLength: 0\r\n\r\n")
			return
		}
		req := Request{
			Method: method,
			Path:   path,
			Body:   body,
		}
		PostFile(conn, req)
	} else {
		fmt.Fprintf(conn, "STATUS 400\r\nLength: 0\r\n\r\n")
	}
}
func GetFile(conn net.Conn, req Request) {
	clean := strings.TrimPrefix(req.Path, "/")
	final := "./" + clean
	db, err := os.ReadFile(final)
	if err != nil {
		errMsg := "Resource not found"
		fmt.Fprintf(conn, "STATUS 404\r\nLength: %d\r\n\r\n%s", len(errMsg), errMsg)
		return
	}
	fmt.Fprintf(conn, "STATUS 200\r\nLength: %d\r\n\r\n", len(db))
	_, err = conn.Write(db)
	if err != nil {
		fmt.Println("Error writing file data:", err)
	}
}

func PostFile(conn net.Conn, req Request) {
	clean := strings.TrimPrefix(req.Path, "/")
	final := "./" + clean
	err := os.WriteFile(final, req.Body, 0644)
	if err != nil {
		fmt.Fprintf(conn, "STATUS 500\r\nLength: 0\r\n\r\n")
	} else {
		fmt.Fprintf(conn, "STATUS 200\r\nLength: 0\r\n\r\n")
	}
}
