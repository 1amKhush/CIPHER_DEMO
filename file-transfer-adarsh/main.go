package main



import (	
	"fmt"
	"net/http"
	"file-transfer-adarsh/controllers"
)

/*

fr -- file should be sent through the network 
	file selector  sender end 
	check for the correct method of seding the file 
	connected device to send the file 
	reciever end. 




*/
func main() {
	fmt.Println("Starting Cipher File Transfer Web UI on http://localhost:3000 🌐")
		
	// Setup API routes
	http.HandleFunc("/api/upload", controllers.HandleUpload)
	http.HandleFunc("/api/download", controllers.HandleDownload)
	http.HandleFunc("/api/start-receiver", controllers.HandleStartReceiver)

	// Serve static frontend
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting web server:", err)
	}
}



