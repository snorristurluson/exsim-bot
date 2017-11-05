package main

import (
	"net"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"math/rand"
	"time"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type SimpleBot struct {
	userid int64
	connection net.Conn
	targetLocation Vector3
}

func NewSimpleBot(userid int64) (*SimpleBot) {
	return &SimpleBot{
		userid: userid,
	}
}

func (bot* SimpleBot) Connect(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Couldn't connect to %v", address)
		return
	}
	fmt.Printf("Connection established\n")
	bot.connection = conn

	msg := fmt.Sprintf(`{"user": %d}` + "\n", bot.userid)
	bot.connection.Write([]byte(msg))
}

func (bot* SimpleBot) SetTargetLocation(loc Vector3) {
	msg := fmt.Sprintf(`{"settargetlocation": {"x": %v, "y": %v, "z": %v}}` + "\n", loc.X, loc.Y, loc.Z)
	bot.connection.Write([]byte(msg))

}

func (bot* SimpleBot) ReceiveLoop() {
	for {
		recvBuf := make([]byte, 4096)
		_, err := bot.connection.Read(recvBuf)
		if err != nil {
			fmt.Printf("Error in Read: %v\n", err)
			break
		}
	}
}

func (bot* SimpleBot) BehaviorLoop() {
	for bot.connection != nil {
		bot.SetTargetLocation(Vector3{X: rand.Float64() * 800, Y: rand.Float64() * 600, Z: 0.0})
		time.Sleep(10 * time.Second)
	}
}

func main() {
	for i:= 0; i < 100; i++ {
		bot := NewSimpleBot(int64(i))
		bot.Connect("localhost:4040")
		if bot.connection == nil {
			return
		}
		go bot.ReceiveLoop()
		go bot.BehaviorLoop()

		time.Sleep(1 * time.Second)
	}

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
