package main

import(
	"github.com/gorilla/mux"
	"net/http"
	
)
func Router () *mux.Router{
	router:=mux.NewRouter()
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/api/user", user).Methods("POST")
	router.HandleFunc("/api/login", login).Methods("POST")

	//router.Use(Authmiddleware)-can be used to protect routes after it
    // router.HandleFunc("/api/notes", GetAllNotes).Methods("GET")
    // router.HandleFunc("/api/note", CreateNote).Methods("POST")
    // router.HandleFunc("/api/note/{id}", UpdateNote).Methods("PUT")
    // router.HandleFunc("/api/note/{id}", DeleteNote).Methods("DELETE")
    // router.HandleFunc("/api/notes", DeleteALLNotes).Methods("DELETE")
	// return router

	// other way:-
    protected:= router.PathPrefix("/api").Subrouter() //PathPrefix is used to group routes under a common prefix
	protected.Use(Authmiddleware)
    
	protected.HandleFunc("/api/notes", GetAllNotes).Methods("GET")
    protected.HandleFunc("/api/note", CreateNote).Methods("POST")
    protected.HandleFunc("/api/note/{id}", UpdateNote).Methods("PUT")
    protected.HandleFunc("/api/note/{id}", DeleteNote).Methods("DELETE")
    protected.HandleFunc("/api/notes", DeleteALLNotes).Methods("DELETE")
	return router 
}

func home(w http.ResponseWriter,r*http.Request){
	w.Write([]byte("Welcome to the home page!"))
}