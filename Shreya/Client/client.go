package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	PostJsonRequest()
}

func PostJsonRequest() {

	requestBody := strings.NewReader(`{"text": "Hi!My name is Shreya"}`)

	response, err := http.Post("http://localhost:9999/message", "application/json", requestBody)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	content, _ := io.ReadAll(response.Body)
	var responseString strings.Builder
	responseString.Write(content)
	fmt.Println(responseString.String())
}
