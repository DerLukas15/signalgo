// Package signalgo provides a dbus interface to talk to the signal-cli from https://github.com/AsamK/signal-cli.
//
// Using DBus should be preferred as it guarantees fixed formating and type definitions.
// Furthermore, more processes can use the signal-cli simultaneously when using the DBus implementation.
// Systembus or Sessionbus are both possible to choose from during the creation.
// Message handling and other events are managed through the Event interface.
//
// Copyright 2021 Lukas Gallandi. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.
package signalgo

import (
	"sync"

	dbus "github.com/godbus/dbus/v5"
	"github.com/sirupsen/logrus"
)

// Event defines all the possible Methods which can be called depending on the received DBus events.
type Event interface {
	OnMessage(source, message string, attachments []string, messagetimestamp int64) error
	OnMessageRead(source string, messagetimestamp int64) error
}

// Connection is the main struct which stores necessary information and is used for
// controlling and using the signalMessenger.
type Connection struct {
	connection   *dbus.Conn
	connectionOK bool

	eventHandler  Event
	eventListener chan bool

	running bool
	replies map[uint32]chan *dbus.Message

	logger *logrus.Entry

	lock sync.RWMutex
}

// New creates a new connection object and starts the necessary routines.
// If no event handling is desired, nil can be used safely.
func New(useSystembus bool, eventHandler Event) (conn *Connection, err error) {
	conn = &Connection{
		logger: logrus.WithFields(logrus.Fields{
			"Routine": "signal",
		}),
		replies: make(map[uint32]chan *dbus.Message),
	}
	if useSystembus {
		conn.connection, err = dbus.SystemBus()
	} else {
		conn.connection, err = dbus.SessionBus()
	}
	if err == nil {
		conn.connectionOK = true
	}
	if eventHandler != nil {
		conn.eventHandler = eventHandler
	}
	go conn.run()
	return
}

// SetEventHandler sets a new eventHandler for the connection.
func (conn *Connection) SetEventHandler(eventHandler Event) (err error) {
	conn.eventHandler = eventHandler
	return
}

// SetLogger sets a new logger for the connection.
func (conn *Connection) SetLogger(logger *logrus.Entry) (err error) {
	conn.logger = logger.WithFields(logrus.Fields{
		"Routine": "signal",
	})
	return
}

// Close closes the connection to the DBus and stopes all running routines.
func (conn *Connection) Close() {
	if conn.running {
		if conn.eventListener != nil {
			close(conn.eventListener)
		}
	}
	if conn.connection != nil {
		conn.connection.Close()
	}
}
