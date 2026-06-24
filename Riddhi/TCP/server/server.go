package main
import (
	"net"
	"fmt"
    "io"
	"log"
	"os"
    "bufio"
    "strings"
)
const (
    HOST = "localhost"
    PORT = "9001"
    TYPE = "tcp"
)

func main(){
    listen, err := net.Listen(TYPE, HOST + ":" + PORT)//The Listen function creates servers
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
    println("Received a request: " + conn.RemoteAddr().String())
    reader:=bufio.NewReader(conn)
    method,err:= reader.ReadString(' ')
    if err != nil {
        log.Fatal(err)
    }
    method = strings.TrimSpace(method)
    filename,err:= reader.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }
    filename = strings.TrimSpace(filename)
    switch method {
    case "POST":
          handleUpload(conn, reader, filename)
    case "GET":
         handleDownload(conn,filename)
       default:
        Status := "STATUS 400\n\nBad Request"
        conn.Write([]byte(Status))
            log.Fatal("Unsupported method")    
     }
     defer conn.Close()
}

func handleUpload(conn net.Conn, reader *bufio.Reader, filename string) {
    fmt.Printf("Handling upload for file: %s\n", filename)
    contentLength, err := reader.ReadString('\n')
    reader.ReadString('\n')
    var size int
    fmt.Sscanf(contentLength, "Content-Length: %d", &size)
    data:= make([]byte,size)
    _, err = io.ReadFull(reader, data)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Received file data of size: %d bytes\n", size)
    fmt.Printf("Data Received = %q\n", string(data))
    // Process the uploaded data and save it to a file
    path := "storage/" + filename
    fmt.Printf("Saving uploaded file to: %s\n", path)
     if _, err := os.Stat(path); err == nil {
    conn.Write([]byte("STATUS 409\n\nFile Already Exists"))
       return
     }
    file, err := os.Create(path)
    if err != nil {
        Status := "STATUS 500\n\nInternal Server Error\n"
        conn.Write([]byte(Status))
        log.Fatal(err)
    }
    defer file.Close()
    _, err = file.Write(data)
       if err != nil {
         log.Fatal(err)
          }
    fmt.Println("FILE WRITTEN in: " + path)
    response := "STATUS 200\n\nUploaded Successfully\n"
    conn.Write([]byte(response))
}

func handleDownload(conn net.Conn,filename string){
    path := "storage/" + filename
    file, err := os.Open(path)
    if err != nil {
        Status := "STATUS 500\nInternal Server Error\n"
        conn.Write([]byte(Status))
        log.Fatal(err)
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        log.Fatal(err)
    }

    size := fileInfo.Size()
    data := make([]byte, size)
    _, err = io.ReadFull(file, data)
    if err != nil {
        log.Fatal(err)
    }

    responseStr := fmt.Sprintf(
        "STATUS %d Downloaded successfully\nContent-Length: %d\n\n",
        200,
        uint32(size),
    )

    conn.Write([]byte(responseStr))
    conn.Write(data)
}

// func handleIncomingRequests(conn net.Conn){
//     println("Received a request: " + conn.RemoteAddr().String());
//     headerBuffer := make([]byte, 1024);

// _, err := io.ReadFull(conn, headerBuffer)
//     if err != nil {
//         log.Fatal(err);
//     }

//     var name string;
//     var reps uint32;

//     if(headerBuffer[0] == byte(1) && headerBuffer[1023] == byte(0)){
//         reps = binary.BigEndian.Uint32(headerBuffer[1:5]);
//         lengthOfName := binary.BigEndian.Uint32(headerBuffer[5:9]);
//         name = string(headerBuffer[9:9+lengthOfName]);
//     } else {
//         log.Fatal("Invalid header");
//     }

//     conn.Write([]byte("Header Received"));

//     dataBuffer := make([]byte, 1024);

//     name = filepath.Base(name);
//     file, err := os.Create(name);
//     if err != nil {
//         log.Fatal(err);
//     }

//     for i := 0; i<int(reps); i++ {
//         _, err := io.ReadFull(conn, dataBuffer)
//         if err != nil {
//             log.Fatal(err);
//         }

//         if(dataBuffer[0] == byte(0) && dataBuffer[1023] == byte(1)){
//             segmentNumber := dataBuffer[1:5];
//             fmt.Printf("Segment Number: %d\n", binary.BigEndian.Uint32(segmentNumber));
//             length := binary.BigEndian.Uint32(dataBuffer[5:9]);
//             fmt.Printf("File Data: %s\n", hex.EncodeToString(dataBuffer[9:9+length]));
//             file.Write(dataBuffer[9:9+length]);
//         } else {
//             log.Fatal("Invalid Segment");
//         }

//          conn.Write([]byte("Segment Received\n"));
//     }

//     time := time.Now().UTC().Format("Monday, 02-Jan-06 15:04:05 MST");
//     conn.Write([]byte("TRANSFER COMPLETE: " + time))

//     file.Close();
//     conn.Close();
// }
