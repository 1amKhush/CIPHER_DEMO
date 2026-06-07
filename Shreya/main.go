package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Message struct {
	Text string `json:"text"`
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	databytes, _ := io.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(databytes)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/message", getMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":9999", router))

}
