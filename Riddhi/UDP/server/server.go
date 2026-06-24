package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "strings"

)

func main() {
    // Set up UDP address
    addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
    if err != nil {
        log.Fatal("Couldn't resolve address:", err)
    }

    // Start listening
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Fatal("Listen failed:", err)
    }
    defer conn.Close()
	println("Server has started")

    for {
        buffer := make([]byte, 1024)
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            log.Printf("Read error: %v", err)
            continue
        }
        fmt.Printf("Got message from %s:\n", clientAddr)
        fmt.Printf("%s\n", string(buffer[:n]))
        go handleRequest(conn,clientAddr, string(buffer[:n]))
    }
}

func handleRequest( conn *net.UDPConn,clientAddr *net.UDPAddr,request string,){
    lines := strings.Split(request, "\n")
    firstLine := strings.Fields(lines[0])
    method := firstLine[0]
    filename := firstLine[1]
    method = strings.TrimSpace(method)
    filename = strings.TrimSpace(filename)
    switch method {
    case "POST":
          handleUpload(conn, clientAddr, filename,request)
    case "GET":
         handleDownload(conn, clientAddr, filename)
       default:
        Status := "STATUS 400\n\nBad Request"
        conn.WriteToUDP([]byte(Status), clientAddr)
            log.Fatal("Unsupported method")    
     }
}

func handleUpload(conn *net.UDPConn, clientAddr *net.UDPAddr, filename string, request string) {
    lines := strings.Split(request, "\n")
    dataLine := lines[2]
    data := strings.TrimPrefix(dataLine, "data:")
    // Receive file data from client
    size := len(data)
    fmt.Printf("Handling upload for file: %s\n", filename)
    fmt.Printf("Received file data of size: %d bytes\n", size)
    fmt.Printf("Data Received = %q\n", string(data))
    // Process the uploaded data and save it to a file
    path := "storage/" + filename
    fmt.Printf("Saving uploaded file to: %s\n", path)
     if _, err := os.Stat(path); err == nil {
    conn.WriteToUDP([]byte("STATUS 409\n\nFile Already Exists"), clientAddr)
       return
     }
    file, err := os.Create(path)
    if err != nil {
        Status := "STATUS 500\nInternal Server Error\n"
        conn.WriteToUDP([]byte(Status), clientAddr)
        log.Fatal(err)
    }
    defer file.Close()
    _, err = file.Write([]byte(data))
       if err != nil {
         log.Fatal(err)
          }
     fmt.Println("FILE WRITTEN in: " + path)
    response := "STATUS 200\n\nUploaded Successfully\n"
    conn.WriteToUDP([]byte(response), clientAddr)
}

func handleDownload(conn *net.UDPConn, clientAddr *net.UDPAddr, filename string){
    path := "storage/" + filename
    file, err := os.Open(path)
    if err != nil {
        Status := "STATUS 500\nInternal Server Error\n"
        conn.WriteToUDP([]byte(Status), clientAddr)
        log.Fatal(err)
    }
    defer file.Close()
    fileInfo, err := file.Stat()
    if err != nil {
        log.Fatal(err)
    }
    size := fileInfo.Size()
    data := make([]byte, size)
    _, err = file.Read(data)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sending file data of size: %d bytes\n", size)
    fmt.Printf("Data Sent = %q\n", string(data))
    // Send file data to client
    _, err = conn.WriteToUDP(data, clientAddr)
    if err != nil {
        log.Fatal(err)
    }
    
}


    // // Buffer for incoming data
    // buffer := make([]byte, 1024)
    // for {
    //     // Read client message
    //     n, clientAddr, err := conn.ReadFromUDP(buffer)
    //     if err != nil {
    //         log.Printf("Read error: %v", err)
    //         continue
    //     }
    //     fmt.Printf("Got message from %s: %s\n", clientAddr, string(buffer[:n]))

    //     // Echo back
    //     _, err = conn.WriteToUDP(buffer[:n], clientAddr)
    //     if err != nil {
    //         log.Printf("Write error: %v", err)
    //     }
    // }


