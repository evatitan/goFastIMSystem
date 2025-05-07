package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	Conn   net.Conn
	Server *Server
}

// NewUser creates a new User instance.
func NewUser(conn net.Conn, server *Server) *User {
	// Get the remote address of the connection
	// This is the address of the client that connected to the server
	// The address is in the format "IP:Port"
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Conn:   conn,
		Server: server,
	}

	// create a goroutine to listen for messages
	// This goroutine will listen for messages sent to the user's channel
	go user.ListenMessage()

	return user

}

// Online sends a message to the user indicating that they are online.
func (this *User) Online() {
	// add user in online map when get connection
	this.Server.mapLock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.mapLock.Unlock()
	// publish this "online" message to all users
	this.Server.Broadcast(this, "online")
}

// Offline sends a message to the user indicating that they are offline.
func (this *User) Offline() {
	// delete user in online map when quit connection
	this.Server.mapLock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.mapLock.Unlock()
	// publish this "online" message to all users
	this.Server.Broadcast(this, "offline")
}

// sendMsg to the others users
func (this *User) DoMessage(msg string) {
	// send message to all users
	sendMsg := "[" + this.Addr + "]" + "," + this.Name + ":" + msg
	this.Server.Message <- sendMsg
	this.Server.Broadcast(this, msg)
}

// ListenMessage listens for messages on the user's channel and sends them to the client.
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg + "\n"))
	}
}
