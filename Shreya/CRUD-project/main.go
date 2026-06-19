package main

import (
	"CRUD-project/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/note/create", controller.RequireAuth(controller.CreateNote)).Methods("POST")
	r.HandleFunc("/note/read", controller.RequireAuth(controller.ReadNote)).Methods("GET")
	r.HandleFunc("/note/update", controller.RequireAuth(controller.UpdateNote)).Methods("PUT")
	r.HandleFunc("/note/delete", controller.RequireAuth(controller.DeleteNote)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4000", r))
}
