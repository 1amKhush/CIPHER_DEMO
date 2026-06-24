package main

import (
	"encoding/json"
	"fmt"

	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)
func main(){
	router := mux.NewRouter() 
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/message",message).Methods("POST")
	router.HandleFunc("/json", jsonmessage).Methods("POST")
	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
func home(w http.ResponseWriter,r*http.Request){
	w.Write([]byte("Welcome to the home page!"))
}
func message(w http.ResponseWriter,r*http.Request){
	fmt.Println("Request method:",r.Method)
	body,err:=io.ReadAll(r.Body)
	if err != nil{
		panic(err)
	}
	message:=string(body)
	fmt.Println("Received:",message)
	w.Write([]byte("Server received :"+ message))
}

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}

func jsonmessage(w http.ResponseWriter, r *http.Request) {

	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	jsonData1, _ := json.MarshalIndent(req, "", "  ")
    fmt.Println("Received:", string(jsonData1))

	resp := Response{
		Reply: "Server received: " + req.Message,
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData2, err := json.MarshalIndent(resp, "", "  ")
    if err != nil {
  panic(err)
 }
    w.Write(jsonData2)

	//other way:-
	// err=json.NewEncoder(w).Encode(resp)
	// if err!=nil{
	// 	panic(err)
	// }
}