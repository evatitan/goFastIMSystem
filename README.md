# goFastIMSystem

`# goFastIMSystem` is a lightweight instant messaging (IM) application written in Go. It demonstrates the implementation of a simple client-server architecture for real-time communication. The project supports features like public chat, private messaging, and user management.

## Features

- **Public Chat**: Users can send messages to all connected users in a shared chatroom.
- **Private Messaging**: Users can send direct messages to specific users.
- **User Management**:
  - View online users.
  - Change username dynamically.
- **Real-Time Communication**: Messages are sent and received in real-time using TCP connections.

## Project Structure

- **`server.go`**: Implements the server-side logic, including managing connected users, broadcasting messages, and handling private messages.
- **`client.go`**: Implements the client-side logic, including connecting to the server, sending/receiving messages, and interacting with the user via a terminal-based menu.
- **`README.md`**: Documentation for the project.

## How It Works

1. **Server**:
   - Listens for incoming TCP connections.
   - Manages a map of online users.
   - Handles public and private messages.
   - Broadcasts messages to all connected users.

2. **Client**:
   - Connects to the server using a TCP connection.
   - Provides a menu for users to interact with the system.
   - Sends commands and messages to the server.
   - Displays messages received from the server in real-time.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/goFastIMSystem.git
   cd goFastIMSystem

2. Build the server and client:
  go build server.go
  go build client.go

## Usage
1. Start the Server
Run the server to start listening for client connections:

2. Start the Client
Run the client and connect to the server:

- Replace <server_ip> with the server's IP address (default: 127.0.0.1).
- Replace <server_port> with the server's port (default: 8888).
- Client Menu, once connected, the client displays the following menu:

1: Send a message to all users in the public chat.
2: Start a private chat with a specific user.
3: Change your username.
0: Disconnect from the server.

3. Example
- Start the server:
  ```bash
  ./server
- Start two clients and connect them to the server:
  ```bash
  ./client -ip 127.0.0.1 -port 8888
- Use the menu to send public and private messages between the clients.  