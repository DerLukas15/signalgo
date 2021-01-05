package main

import (
	"errors"
	"fmt"

	"github.com/DerLukas15/signalgo"
)

type eventHandler struct{}

func (eH *eventHandler) OnMessage(source, message string, attachments []string, messagetimestamp int64) error {
	fmt.Println("Got Message ", message, " (", messagetimestamp, ") from ", source, " with attachments ", attachments)
	return nil
}

func (eH *eventHandler) OnMessageRead(source string, messagetimestamp int64) error {
	fmt.Println("Message from ", source, " with timestamp ", messagetimestamp, " has been received")
	return nil
}

func main() {
	eventHandler := &eventHandler{}
	conn, err := signalMessenger.New(true, eventHandler)
	if err != nil {
		panic(err)
	}
	if conn == nil {
		panic(errors.New("No connection"))
	}
	for {
	}
}
