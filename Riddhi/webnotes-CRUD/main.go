package main
import(
	"fmt"
	"net/http"
	"log"
	"github.com/joho/godotenv"
	
)


func main(){
	//if you are using .env file to store your environment variables, load them using godotenv
	err := godotenv.Load()  
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	ConnectDB()
	r:=Router()
	fmt.Println("Server running on :8081")
    log.Fatal(http.ListenAndServe(":8081", r))

}


