package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	PORT = flag.Int("port", 35000, "Listening port")
)

var clients = []net.Conn{}

func connect() {
	address := fmt.Sprintf(":%d", *PORT)
	log.Printf("Trying to listen on %s...\n", address)
	conn, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Failed to listen on %s\n", address)
	}
	defer conn.Close()
	log.Printf("Successfully listening on %s\n", address)

	for {
		client, err := conn.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s\n", err)
			return
		}
		// log.Printf("Accepted connection from %s\n", client.RemoteAddr())
		go handle(client)
	}
}

func handle(client net.Conn) {
	defer func() {

		ret := make([]net.Conn, 0, len(clients))
		for _, val := range clients {
			if val != client {
				ret = append(ret, val)
			}
		}
		clients = ret

		client.Close()
	}()

	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil {
		log.Printf("Failed to read from client: %s\n", err)
		return
	}
	str := trimString(string(buf[:n]))

	if !strings.HasPrefix(str, "Username: ") {
		_, _ = client.Write([]byte("Invalid username"))
		return
	}

	username := strings.TrimPrefix(str, "Username: ")
	username = trimString(username)

	if username == "" {
		_, _ = client.Write([]byte("Username is empty"))
		return
	}

	clients = append(clients, client)

	welcomeMsg := fmt.Sprintf("%s has joined the chat", username)

	log.Println(welcomeMsg)

	broadcast(client, fmt.Sprintf("[Server]: %s\n", welcomeMsg))

	for {
		n, err = client.Read(buf)

		if err != nil {
			break
		}

		str = trimString(string(buf[:n]))

		broadcast(client, fmt.Sprintf("%s: %s\n", username, str))
	}
}

func broadcast(client net.Conn, msg string) {
	for _, c := range clients {
		if c != client {
			_, _ = c.Write([]byte(msg))
		}
	}
}

func trimString(str string) string {

	str = strings.Trim(str, "\r\n")
	str = strings.Trim(str, "\n")
	str = strings.TrimSpace(str)

	return str
}

func main() {
	flag.Parse()
	if *PORT < 1 || *PORT > 65535 {
		log.Fatal("Port is invalid")
	}
	for {
		connect()
	}
}
