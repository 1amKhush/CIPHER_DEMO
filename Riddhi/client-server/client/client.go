package main 
import(
	"fmt"
	"io"
	"net/http"
	"strings"
	"encoding/json"
	"bytes"
)
type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}
func main(){
	message:="Hello, Server!"
	
	textresp,err:=http.Post("http://localhost:8081/message","text/plain",strings.NewReader(message))
	if err != nil {
		panic(err)
	}
	defer textresp.Body.Close()
	textbody, err := io.ReadAll(textresp.Body) 
	if err != nil {
		panic(err)
	}	
	fmt.Println("Response:")
	fmt.Println(string(textbody))
	fmt.Println(textresp.Status)

	req := Request{
	Message: "Hello, Server!",
    }
    jsonData, err := json.MarshalIndent(req, "", " ")
    if err != nil {
       panic(err)
    }

     jsonreq,err:=http.Post("http://localhost:8081/json","application/json",bytes.NewBuffer(jsonData))

     jsonbody, _ := io.ReadAll(jsonreq.Body)
     fmt.Println(string(jsonbody))
	 
	 defer jsonreq.Body.Close()
}