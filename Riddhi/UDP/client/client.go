package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "io"
)

func main() {
    // Connect to server
    addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
    if err != nil {
        log.Fatal("Couldn't resolve address:", err)
    }

    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        log.Fatal("Connection failed:", err)
    }
    defer conn.Close()

    if len(os.Args) < 3 {
        log.Fatal("Usage: go run client.go <POST/GET> <filename>")
    }
    
    switch method := os.Args[1]; method {
    case "POST":
        file, err := os.Open(os.Args[2])
        if err != nil {
            log.Fatal(err)
        }
        uploadFile(conn, file)
        response := make([]byte, 1024)
        n, err := conn.Read(response)
        if err != nil {
        log.Fatal(err)
        }
       fmt.Println(string(response[:n]))
    case "GET":
        downloadFile(conn, os.Args[2])
        response := make([]byte, 1024)
        n, err := conn.Read(response)
        if err != nil {
        log.Fatal(err)
        }
       fmt.Println(string(response[:n]))
    default:
        log.Fatal("Invalid method")
    }
}

type Request struct {
    method string
    filenamelength uint32
    filename string
    segmentSize uint32 // segment size for file transfer
    reps uint32  // number of times to repeat the file transfer for testing
    dataSize uint32
    data []byte
}

func uploadFile(conn *net.UDPConn, file *os.File) {

      fileInfo, err := file.Stat()
      if err != nil {
        log.Fatal(err)
       }
         size := fileInfo.Size()
         name := fileInfo.Name()
    
     data, err := io.ReadAll(file)//Difference between ReadAll and ReadAt is that ReadAll reads the entire file into memory, while ReadAt reads a specific portion of the file based on the provided offset and length. ReadAll is simpler to use when you want to read the entire file, while ReadAt is more efficient for reading specific parts of a large file without loading it entirely into memory.
      if err != nil {
     log.Fatal(err)
         }

    request:= Request{
        method: "POST",
        filenamelength: uint32(len(name)),
        filename: name,
        segmentSize: 1024, // example segment size
        reps: 1, // example repetition count
        dataSize: uint32(size),
        data: data,
    }

    requestStr := fmt.Sprintf(
    "POST %s\nContent-Length: %d\ndata:%s",
    request.filename,
    request.dataSize,
    string(request.data),
      )

       conn.Write([]byte(requestStr))
       
}

func downloadFile(conn *net.UDPConn, filename string) {
    requestStr := fmt.Sprintf(
        "GET %s\n\n",
        filename,
    )
    conn.Write([]byte(requestStr))
}

 // // Send a message
    // message := []byte("Hello, UDP!")
    // _, err = conn.Write(message)
    // if err != nil {
    //     log.Printf("Send failed: %v", err)
    //     return
    // }

    // // Wait for reply with timeout
    // conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    // buffer := make([]byte, 1024)
    // n, _, err := conn.ReadFromUDP(buffer)
    // if err != nil {
    //     log.Printf("Receive error: %v", err)
    //     return
    // }
    // fmt.Printf("Server says: %s\n", string(buffer[:n]))