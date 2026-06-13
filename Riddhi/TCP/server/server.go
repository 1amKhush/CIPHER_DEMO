package main
import (
	"net"
	"encoding/binary"
	"fmt"
    "io"
	"log"
	"os"
    "path/filepath"
    "encoding/hex"
    "time"
)
const (
    HOST = "localhost"
    PORT = "9001"
    TYPE = "tcp"
)

func main(){
    listen, err := net.Listen(TYPE, HOST + ":" + PORT)
    if err != nil {
        log.Fatal(err)
    }
   
    defer listen.Close()

    println("Server has started on PORT " + PORT)

    for {
        conn, err := listen.Accept()
        if err != nil {
            log.Fatal(err)
        }
        println("Hello");
        go handleIncomingRequests(conn)
    }
}
func handleIncomingRequests(conn net.Conn){
    println("Received a request: " + conn.RemoteAddr().String());
    headerBuffer := make([]byte, 1024);

_, err := io.ReadFull(conn, headerBuffer)
    if err != nil {
        log.Fatal(err);
    }

    var name string;
    var reps uint32;

    if(headerBuffer[0] == byte(1) && headerBuffer[1023] == byte(0)){
        reps = binary.BigEndian.Uint32(headerBuffer[1:5]);
        lengthOfName := binary.BigEndian.Uint32(headerBuffer[5:9]);
        name = string(headerBuffer[9:9+lengthOfName]);
    } else {
        log.Fatal("Invalid header");
    }

    conn.Write([]byte("Header Received"));

    dataBuffer := make([]byte, 1024);

    name = filepath.Base(name);
    file, err := os.Create(name);
    if err != nil {
        log.Fatal(err);
    }

    for i := 0; i<int(reps); i++ {
        _, err := io.ReadFull(conn, dataBuffer)
        if err != nil {
            log.Fatal(err);
        }

        if(dataBuffer[0] == byte(0) && dataBuffer[1023] == byte(1)){
            segmentNumber := dataBuffer[1:5];
            fmt.Printf("Segment Number: %d\n", binary.BigEndian.Uint32(segmentNumber));
            length := binary.BigEndian.Uint32(dataBuffer[5:9]);
            fmt.Printf("File Data: %s\n", hex.EncodeToString(dataBuffer[9:9+length]));
            file.Write(dataBuffer[9:9+length]);
        } else {
            log.Fatal("Invalid Segment");
        }

         conn.Write([]byte("Segment Received\n"));
    }

    time := time.Now().UTC().Format("Monday, 02-Jan-06 15:04:05 MST");
    conn.Write([]byte("TRANSFER COMPLETE: " + time))

    file.Close();
    conn.Close();
}
