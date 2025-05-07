package main

import (
	"fmt"
	"net"
	"strings"
)

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

// sendMsg sends a message to this user pc

func (this *User) SendMsg(msg string) {
	_, err := this.Conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("Error sending message to", this.Name, ":", err)
	}
}

// sendMsg to the others users
func (this *User) DoMessage(msg string) {
	// search all online users
	if msg == "who" {
		this.Server.mapLock.Lock()
		for _, user := range this.Server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + "," + user.Name + ":" + "online\n"
			this.SendMsg(onlineMsg)
		}
		this.Server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// get the new name
		newName := strings.Split(msg, "|")[1]
		// check if the name is already used
		if _, ok := this.Server.OnlineMap[newName]; ok {
			this.SendMsg("this name is already used\n")
		} else {
			// rename the user
			this.Server.mapLock.Lock()
			delete(this.Server.OnlineMap, this.Name)
			// update the key in the map
			/*/ OnlineMap = map[string]*User{
			    "Alice": &User{Name: "Alice", Addr: "127.0.0.1"},
			    "Bob":   &User{Name: "Bob", Addr: "127.0.0.2"},
			/ }*/
			this.Server.OnlineMap[newName] = this
			// update the value in the map
			this.Name = newName
			this.Server.mapLock.Unlock()
			this.SendMsg("You renamed to " + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// msg: "to|Bob|hello"
		// get msg receiver name
		parts := strings.Split(msg, "|")
		if len(parts) < 3 {
			this.SendMsg("Invalid message format. Use: to|<receiver_name>|<message_content>\n")
			return
		}
		remoteName := parts[1]
		msgContent := parts[2]

		if remoteName == "" {
			this.SendMsg("Please input valid receiver name\n")
			return
		}
		// check if the receiver is online
		remoteUser, ok := this.Server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("the receiver is not online \n")
			return
		}

		// get the msg content
		if msgContent == "" {
			this.SendMsg("Please input valid content\n")
			return
		}

		// send the msg to the receiver
		remoteUser.SendMsg(this.Name + " sent you a message: " + msgContent + "\n")

	} else {
		// send message to all users
		this.Server.Broadcast(this, msg)
	}

}

// ListenMessage listens for messages on the user's channel and sends them to the client.
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg + "\n"))
	}
}
