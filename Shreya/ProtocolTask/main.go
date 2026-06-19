package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	go s.acceptLoop()
	<-s.quitch
	return nil
}
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("couldn't accept because:", err)
			continue
		}
		go s.Handleconnection(conn)
	}
}
func main() {
	addr := ":3000"
	server := NewServer(addr)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
