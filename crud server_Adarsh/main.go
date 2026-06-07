package main

import (
	"fmt"
	"log"
	"net/http"
)
/*
User Authentication and Authorization :
user get login through certain email and password and get acess to the application 
to prevent the user 
how to store the jwt tokens 
user email and user_name and password 
register then user will send the jwt tokens 
just email and password 

encrypt 

jwkey  // store in env 
struct jwt payload 

database 
{ 
	sql or psql
	in memory database currently
}
	routers 
	handlers 
	register 
	login {
		verify the hash 
		expiration time 
		create the jwt token
		generate the sign 
		return token 

	}
	middleware -- extract and validate the jwt 
	helper context wrapper 



	



*/

// Actually not required in this codebase

func enableCORS(next http.Handler) http.Handler {  
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers globally
		w.Header().Set("Access-Control-Allow-Origin", "*") // Change * to specific origin for security
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests immediately
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}






func main(){
	// Router.HandleFunc("/", ServeHome).Methods("GET")
	fmt.Println("Hello World")

	InitDB()

	r := Route()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":4000", enableCORS(r)))
	fmt.Println("Listening at port 4000 ...")
} 