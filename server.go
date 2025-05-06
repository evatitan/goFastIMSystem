package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer creates a new server instance with the given IP and port.
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   "127.0.0.1",
		Port: 8888,
	}
	return server
}

func (this *Server) handler(conn net.Conn) {
	fmt.Println("connection accepted")
}

func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// socket close
	defer listener.Close()

	for {

		// socket accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		// handle connection
		go this.handler(conn)
	}

}
