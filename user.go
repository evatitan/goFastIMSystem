package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
}

// NewUser creates a new User instance.
func NewUser(conn net.Conn) *User {
	// Get the remote address of the connection
	// This is the address of the client that connected to the server
	// The address is in the format "IP:Port"
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}

	// create a goroutine to listen for messages
	// This goroutine will listen for messages sent to the user's channel
	go user.ListenMessage()

	return user

}

// ListenMessage listens for messages on the user's channel and sends them to the client.
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg + "\n"))
	}
}
