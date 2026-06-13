package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // Set up UDP address
    addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
    if err != nil {
        log.Fatal("Couldn’t resolve address:", err)
    }

    // Start listening
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Fatal("Listen failed:", err)
    }
    defer conn.Close()
	println("Server has started")

    // Buffer for incoming data
    buffer := make([]byte, 1024)
    for {
        // Read client message
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            log.Printf("Read error: %v", err)
            continue
        }
        fmt.Printf("Got message from %s: %s\n", clientAddr, string(buffer[:n]))

        // Echo back
        _, err = conn.WriteToUDP(buffer[:n], clientAddr)
        if err != nil {
            log.Printf("Write error: %v", err)
        }
    }
}