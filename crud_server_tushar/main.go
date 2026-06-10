package main

import (
	"crud_server/database"
	"fmt"
	"log"
	"net/http"

	"crud_server/routes"
)

func main() {
	fmt.Println("hehe, start of smth newww")

	router := routes.Register_routes()

	database.Connect_db()

	log.Fatal(http.ListenAndServe(":8000", router))
}
