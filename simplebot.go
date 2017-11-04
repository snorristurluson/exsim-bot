package main

import (
	"net"
	"fmt"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4040")
	if err != nil {
		fmt.Errorf("Couldn't connect")
		return
	}

	fmt.Printf("Connection established")

	msg := "{\"user\": 0}\n"
	_, err = conn.Write([]byte(msg))

	if err != nil {
		fmt.Errorf("Couldn't write")
		return
	}

	for {
		recvBuf := make([]byte, 4096)
		n, err := conn.Read(recvBuf)
		if err != nil {
			fmt.Print(err)
			break
		}
		msgReceived := string(recvBuf[:n])
		fmt.Printf("Message Received: %s\n", msgReceived)
	}
}
