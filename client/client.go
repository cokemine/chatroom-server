package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	SERVER   = flag.String("server", "127.0.0.1", "Server address")
	PORT     = flag.Int("port", 35000, "Listening port")
	USERNAME = flag.String("username", "Anonymous", "Username")
)

func connect() {
	address := fmt.Sprintf("%s:%d", *SERVER, *PORT)

	log.Printf("Trying to connect to %s...\n", address)

	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Printf("Failed to connect to %s\n", address)
		return
	}

	defer conn.Close()

	log.Printf("Successfully connected to %s\n", address)

	_, _ = conn.Write([]byte(fmt.Sprintf("Username: %s", *USERNAME)))

	buf := make([]byte, 1024)

	go func() {
		for {
			n, _ := conn.Read(buf)
			fmt.Printf("%s", buf[:n])
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		input, _, _ := reader.ReadLine()
		_, _ = conn.Write(input)
	}
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
