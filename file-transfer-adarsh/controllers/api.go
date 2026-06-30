package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

// HandleUpload receives a file from the Web UI, stages it, and generates a PIN.
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	payload, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Save to staging
	err = os.MkdirAll("staging", 0755)
	if err != nil {
		http.Error(w, "Error creating staging dir", http.StatusInternalServerError)
		return
	}
	filePath := "staging/" + handler.Filename
	err = os.WriteFile(filePath, payload, 0644)
	if err != nil {
		http.Error(w, "Error saving to staging", http.StatusInternalServerError)
		return
	}

	// Generate 6-digit PIN
	pin, err := GeneratePIN(handler.Filename, filePath)
	if err != nil {
		http.Error(w, "Error generating PIN", http.StatusInternalServerError)
		return
	}

	// Send back JSON response with PIN
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"pin":     pin,
		"message": "File uploaded successfully!",
	})
}

// HandleDownload asks for a PIN, validates it, sends it to the TCP Receiver, and downloads it to the browser
func HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pin := r.URL.Query().Get("pin")
	if pin == "" {
		http.Error(w, "PIN required", http.StatusBadRequest)
		return
	}

	// Validate PIN
	transfer, ok := GetTransfer(pin)
	if !ok {
		http.Error(w, "Invalid or expired PIN", http.StatusNotFound)
		return
	}

	// Read from staging
	payload, err := os.ReadFile(transfer.FilePath)
	if err != nil {
		http.Error(w, "Error reading staged file", http.StatusInternalServerError)
		return
	}

	// Construct our custom protocol request to satisfy the assignment
	requestData := FormatRequest("POST", transfer.Filename, payload)

	conn, err := net.Dial("tcp", "localhost:8090")
	if err == nil {
		defer conn.Close()
		_, err = conn.Write(requestData)
		if err == nil {
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				tcpConn.CloseWrite()
			}
			// wait for receiver to process it
			io.ReadAll(conn)
		}
	}
	// Even if the local TCP receiver is down, we still send the file to the browser

	// Clean up staging file
	os.Remove(transfer.FilePath)

	// Send file payload back to user
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", transfer.Filename))
	w.Write(payload)
}

var receiverRunning bool

// HandleStartReceiver starts the TCP receiver in the background
func HandleStartReceiver(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if receiverRunning {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "already running"})
		return
	}

	receiverRunning = true
	// Start receiver in a goroutine so it doesn't block the HTTP response
	go Reciever()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}
