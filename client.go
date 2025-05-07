package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// create a new client instance
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// connect to the server
	conn, err := net.Dial("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Println("net Dial connect to server failed, err:", err)
		return nil
	}

	// return the client instance
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip by default is 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "server port by default is 8888")
}

// client menu
func (client *Client) menu() bool {
	fmt.Println("1. group message")
	fmt.Println("2. private message")
	fmt.Println("3. rename")
	fmt.Println("0. offline")
	var flag int
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("please entre a valid number")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
			fmt.Println("please entre a valid number")
		}
	}
	switch client.flag {
	case 1:
		fmt.Println("group message selected")
		break
	case 2:
		fmt.Println("private message selected")
		break
	case 3:
		fmt.Println("rename selected")
		break
	}
}

func main() {
	// parse command line arguments
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("client connect to server failed")
		return
	}
	fmt.Println("client connect to server success")

	// Run the client menu and handle user input
	client.Run()

	// do other things for client
	select {}
}
