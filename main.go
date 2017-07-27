package main

import (
	"fmt"
	"strconv"
	"time"
)

type Message struct {
	Command string
	Data    string
}
type Event struct {
	Originator string
	Time       string
	Message
}
type iface string

func gendata(i iface) <-chan Message {
	data := make(chan Message)
	counter := 1
	go func() {
		for {
			<-time.After(200 * time.Millisecond)
			msg := Message{
				Command: fmt.Sprintf("cmd : %s", strconv.Itoa(counter)),
				Data:    fmt.Sprintf("data: %s", strconv.Itoa(counter)),
			}
			data <- msg
			counter++
		}
	}()
	return data
}
func NewMon(iface string, c chan struct{}) {
	fmt.Println("Starting new mon", iface)
	c <- struct{}{}
}

func (m Message) String() string {
	return fmt.Sprintf("Message:\n\tCommand: %s\n\tData: %s", m.Command, m.Data)
}
func main() {
	Events := gendata("eth0")
	//interval := time.Minute * 1
	var ifmoncnl chan struct{}

	for {
		select {
		case msg := <-Events:
			fmt.Println(msg.String())
			if msg.Command == "newmon" {
				go NewMon(msg.Data, ifmoncnl)
			}
		default:
		}
		<-time.After(200 * time.Millisecond)
	}
}
