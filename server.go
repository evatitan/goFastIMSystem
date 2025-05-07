package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// online user map
	OnlineMap map[string]*User
	// In Go, maps are NOT concurrency-safe. If two or more goroutines access them simultaneously (read/write), the program will fail with a panic.	mapLock   sync.RWMutex
	mapLock sync.RWMutex
	// message channel
	Message chan string
}

// NewServer creates a new server instance with the given IP and port.
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// MonitorMessage listens for messages on the server's message channel and broadcasts them to all users.
func (this *Server) MonitorMessage() {
	for {
		msg := <-this.Message
		// send message to all users
		this.mapLock.Lock()
		for _, client := range this.OnlineMap {
			client.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Broadcast(user *User, msg string) {
	// send message to all users
	sendMsg := "[" + user.Addr + "]" + ", " + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) handler(conn net.Conn) {
	// handle connection
	// fmt.Println("connection accepted")

	// create a new user
	user := NewUser(conn)

	// add user in online map when get connection
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// publish this "online" message to all users
	this.Broadcast(user, "online")

	// listen and read for user message
	go func() {
		// read message from user
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.Broadcast(user, "offline")
				return
			}
			if err != nil && err.Error() != "EOF" {
				fmt.Println("conn Read err:", err)
				return
			}
			// get message(remove \n)
			// buf[:n] is a slice of the buffer with length n
			msg := string(buf[:n-1])
			// send message to all users
			this.Broadcast(user, msg)
		}
	}()
	//block this goroutine
	select {}

}

func (this *Server) Start() {
	// socket listen ,
	// tcp:Transmission Control Protocol(HTTP web, SSH etc) / UDP: User Datagram Protocol(video, voice, streaming etc)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// socket close
	defer listener.Close()

	// start message monitor
	go this.MonitorMessage()

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
