package main

import (
    "fmt"
    "log"
    "net"
    "time"
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

    // Send a message
    message := []byte("Hello, UDP!")
    _, err = conn.Write(message)
    if err != nil {
        log.Printf("Send failed: %v", err)
        return
    }

    // Wait for reply with timeout
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    buffer := make([]byte, 1024)
    n, _, err := conn.ReadFromUDP(buffer)
    if err != nil {
        log.Printf("Receive error: %v", err)
        return
    }
    fmt.Printf("Server says: %s\n", string(buffer[:n]))
}