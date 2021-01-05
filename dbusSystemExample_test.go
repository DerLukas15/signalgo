// Copyright 2021 Lukas Gallandi. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package signalgo_test

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

func Example_systemDBus() {
	eventHandler := &eventHandler{}
	conn, err := signalgo.NewDBus(true, eventHandler)
	if err != nil {
		panic(err)
	}
	if conn == nil {
		panic(errors.New("No connection"))
	}
	for {
	}
}
