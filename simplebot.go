package main

import (
	"net"
	"fmt"
	"math/rand"
	"time"
	"sync"
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
	bot.connection = conn

	msg := fmt.Sprintf(`{"user": %d}` + "\n", bot.userid)
	bot.connection.Write([]byte(msg))
}

func (bot* SimpleBot) SetTargetLocation(loc Vector3) {
	msg := fmt.Sprintf(
		`{"command": "settargetlocation", "params": {"location": {"x": %v, "y": %v, "z": %v}}}` + "\n",
		loc.X, loc.Y, loc.Z)
	bot.connection.Write([]byte(msg))

}

func (bot* SimpleBot) ReceiveLoop(wg *sync.WaitGroup) {
	for {
		recvBuf := make([]byte, 128*1024)
		_, err := bot.connection.Read(recvBuf)
		if err != nil {
			bot.connection.Close()
			bot.connection = nil
			break
		}
		//s := string(recvBuf[:n])
		//fmt.Printf("Received %v bytes: %s\n", n, s)
	}
	wg.Done()
}

func (bot* SimpleBot) BehaviorLoop() {
	for bot.connection != nil {
		bot.SetTargetLocation(Vector3{X: rand.Float64() * 5000 - 2500, Y: rand.Float64() * 5000 - 2500, Z: 0.0})
		time.Sleep(time.Duration((rand.Float64() * 25 + 5)) * time.Second)
	}
}

func main() {
	numBots := 500
	fmt.Printf("Starting %v bots\n", numBots)

	var wg sync.WaitGroup
	for i:= 0; i < numBots; i++ {
		bot := NewSimpleBot(int64(i))
		bot.Connect("localhost:4040")
		if bot.connection == nil {
			return
		}
		wg.Add(1)
		go bot.ReceiveLoop(&wg)
		go bot.BehaviorLoop()

		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("%v bots running\n", numBots)
	wg.Wait()
}
