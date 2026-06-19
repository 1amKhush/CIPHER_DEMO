package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func TOGET(conn net.Conn, filename string) {
	request := fmt.Sprintf("GET /%s\n", filename)
	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	status = strings.TrimSpace(status)

	var Length int
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading headers:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Length:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				cl, err := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err == nil {
					Length = cl
				}
			}
		}
	}

	if !strings.Contains(status, "STATUS 200") {
		errBody := make([]byte, Length)
		reader.Read(errBody)
		fmt.Printf("Server error: %s\n%s\n", status, string(errBody))
		return
	}

	fileData := make([]byte, Length)
	bytesRead := 0
	for bytesRead < Length {
		n, err := reader.Read(fileData[bytesRead:])
		if err != nil {
			fmt.Println("Error reading file data:", err)
			return
		}
		bytesRead += n
	}

	err = os.WriteFile(filename, fileData, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
	fmt.Printf("Downloaded %s (%d bytes)\n", filename, Length)
}

func TOPOST(conn net.Conn, filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	request := fmt.Sprintf("POST /%s %d\n%s", filename, len(data), string(data))
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	status = strings.TrimSpace(status)

	var Length int
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading headers:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Length:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				cl, err := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err == nil {
					Length = cl
				}
			}
		}
	}
	body := make([]byte, Length)
	if Length > 0 {
		_, err = reader.Read(body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
	}

	if strings.Contains(status, "STATUS 200") {
		fmt.Println("Upload successful")
	} else {
		fmt.Printf("Upload failed: %s\n%s\n", status, string(body))
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("go run client.go GET filename")
		fmt.Println("go run client.go POST filename")
		return
	}

	method := os.Args[1]
	filename := os.Args[2]

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	if method == "GET" {
		TOGET(conn, filename)
	} else if method == "POST" {
		TOPOST(conn, filename)
	} else {
		fmt.Println("Invalid method")
	}
}
