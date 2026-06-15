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

type Response struct {
	Status int
	Size   int
	Body   []byte
}

func ReadResponse(reader *bufio.Reader) (*Response, error) {
	resp := &Response{}

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
			return nil, fmt.Errorf("invalid response")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {

		case "STATUS":

			status, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}

			resp.Status = status

		case "SIZE":

			size, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}

			resp.Size = size
		}
	}

	if resp.Size > 0 {

		resp.Body = make([]byte, resp.Size)

		_, err := io.ReadFull(reader, resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func uploadFile(filename string) error {

	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		return err
	}

	defer conn.Close()

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	request := fmt.Sprintf("METHOD:POST\nPATH:/%s\nSIZE:%d\n\n", filename, len(data))

	_, err = conn.Write([]byte(request))
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	resp, err := ReadResponse(bufio.NewReader(conn))

	if err != nil {
		return err
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println(string(resp.Body))

	return nil
}

func downloadFile(filename string) error {

	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		return err
	}

	defer conn.Close()

	request := fmt.Sprintf("METHOD:GET\nPATH:/%s\n\n", filename)

	_, err = conn.Write([]byte(request))
	if err != nil {
		return err
	}

	resp, err := ReadResponse(bufio.NewReader(conn))

	if err != nil {
		return err
	}

	if resp.Status != 200 {

		fmt.Println(string(resp.Body))
		return nil
	}

	outputFile := filename

	err = os.WriteFile(outputFile, resp.Body, 0644)

	if err != nil {
		return err
	}

	fmt.Println("Downloaded:", outputFile)

	return nil
}

func listFiles() error {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		return err
	}
	defer conn.Close()

	request := "METHOD:LIST\n\n"

	_, err = conn.Write([]byte(request))
	if err != nil {
		return err
	}

	resp, err := ReadResponse(bufio.NewReader(conn))

	if err != nil {
		return err
	}

	fmt.Println("Files on server:")
	fmt.Println(string(resp.Body))

	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("go run main.go upload <file>")
		fmt.Println("go run main.go download <file>")
		fmt.Println("go run main.go list")
		return
	}

	command := os.Args[1]

	switch command {
	case "upload":
		if len(os.Args) < 3 {
			fmt.Println("filename required")
			return
		}

		err := uploadFile(os.Args[2])

		if err != nil {
			fmt.Println(err)
		}

	case "download":
		if len(os.Args) < 3 {
			fmt.Println("filename required")
			return
		}

		err := downloadFile(os.Args[2])

		if err != nil {
			fmt.Println(err)
		}

	case "list":
		err := listFiles()

		if err != nil {
			fmt.Println(err)
		}

	default:
		fmt.Println("Unknown command")
	}
}
