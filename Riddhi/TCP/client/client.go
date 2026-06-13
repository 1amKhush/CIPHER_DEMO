package main
import (
	"net"
	"encoding/binary"
	"fmt"
	"log"
	"os"
    "io"
)
const (
    HOST = "localhost"
    PORT = "9001"
    TYPE = "tcp"
)

type MetaData struct {
    name     string
    fileSize uint32
    reps     uint32
}
func prepareMetadata(file *os.File) MetaData {

    fileInfo, err := file.Stat()

    if err != nil {
        log.Fatal(err)
    }

    size := fileInfo.Size()

    header := MetaData{
        name:     file.Name(),
        fileSize: uint32(size),
        reps:     uint32(size + 1014 - 1) / 1014, // Round up to the nearest multiple of 1014
    }

    return header
}
func main() {
    tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
    if err != nil {
        log.Fatal(err)
    }

    conn, err := net.DialTCP(TYPE, nil, tcpServer)
    if err != nil {
        log.Fatal(err)
    }

    sendFile("file.txt", conn)

}
func sendFile(path string, conn *net.TCPConn) {
    file, err := os.OpenFile(path, os.O_RDONLY, 0755)

    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

    header := prepareMetadata(file)

    dataBuffer := make([]byte, 1014)

    // Start (all 1s) - 1 byte, reps - 4 bytes, lengthofname - 4 bytes, name - `lengthofname` bytes, End (all 0s) - 1 byte;
    headerBuffer := []byte{1}

    // Start (all 0s) - 1 byte, Segment number - 4 bytes, lengthofdata - 4 bytes, Data - `lengthofdata` bytes, End (all 1s) - 1 byte
    segmentBuffer := []byte{0}

    // Temporary buffer for uint32
    temp := make([]byte, 4)

    // Temporary buffer for responses received
    received := make([]byte, 100);

    for i := 0; i < int(header.reps); i++ {
        n,err:= file.ReadAt(dataBuffer, int64(i*1014))
        if err != nil && err != io.EOF {
            log.Fatal(err)
        }

        if i == 0 {
            // Number of segments
            binary.BigEndian.PutUint32(temp, header.reps)
            headerBuffer = append(headerBuffer, temp...)

            // Length of name
            binary.BigEndian.PutUint32(temp, uint32(len(header.name)))
            headerBuffer = append(headerBuffer, temp...)

            // Name
            headerBuffer = append(headerBuffer, []byte(header.name)...)

            headerBuffer = append(headerBuffer, 0)

            for len(headerBuffer) < 1024 {
            headerBuffer = append(headerBuffer, 0)
          }

            _, err := conn.Write(headerBuffer)

            if err != nil {
                log.Fatal(err)
            }

            m, err := conn.Read(received)
            if err != nil {
            log.Fatal(err)
              }

           fmt.Println(string(received[:m]))
        }

        // Segment number
        binary.BigEndian.PutUint32(temp, uint32(i))
        segmentBuffer = append(segmentBuffer, temp...);

        // Length of data
        binary.BigEndian.PutUint32(temp, uint32(n))
        segmentBuffer = append(segmentBuffer, temp...)

        // Data
        segmentBuffer = append(segmentBuffer, dataBuffer...) 
        segmentBuffer = append(segmentBuffer, 1)

        _, err = conn.Write(segmentBuffer);

        if err != nil {
            log.Fatal(err);
        }

       m, err := conn.Read(received)

         if err != nil {
        log.Fatal(err)
         }

       fmt.Println("Server:", string(received[:m]))

        if err != nil {
            log.Fatal(err)
        }

        // Reset segment buffer
        segmentBuffer = []byte{0};
    }
    

}
