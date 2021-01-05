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
	"errors"
	"sync"

	dbus "github.com/godbus/dbus/v5"
	"github.com/sirupsen/logrus"
)

// Event defines all the possible Methods which can be called depending on the received DBus events.
type Event interface {
	OnMessage(source, message string, attachments []string, messagetimestamp int64) error
	OnMessageRead(source string, messagetimestamp int64) error
}

type connectionTyp int64

const (
	connectionDBus connectionTyp = iota + 1
	connectionCLI
)

// Connection is the main struct which stores necessary information and is used for
// controlling and using the signalMessenger.
type Connection struct {
	connection     *dbus.Conn
	executablePath string

	connectionOK   bool
	connectionType connectionTyp

	eventHandler  Event
	eventListener chan bool

	running bool
	replies map[uint32]chan *dbus.Message

	logger *logrus.Entry

	lock sync.RWMutex
}

// NewDBus creates a new connection object which connects to signal-cli via DBus
// and starts the necessary routines.
func NewDBus(useSystembus bool, eventHandler Event) (conn *Connection, err error) {
	conn = &Connection{
		logger: logrus.WithFields(logrus.Fields{
			"Routine": "signalViaDBus",
		}),
		replies:        make(map[uint32]chan *dbus.Message),
		connectionType: connectionDBus,
		eventListener:  make(chan bool),
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
	go conn.runViaDBus()
	return
}

// NewCLI creates a new connection object which connects to signal-cli via CLI
// and starts the necessary routines. The executablePath should be absolute.
func NewCLI(executablePath string, eventHandler Event) (conn *Connection, err error) {
	conn = &Connection{
		logger: logrus.WithFields(logrus.Fields{
			"Routine": "signalViaCLI",
		}),
		replies:        make(map[uint32]chan *dbus.Message),
		connectionType: connectionCLI,
		executablePath: executablePath,
		eventListener:  make(chan bool),
	}

	valid, err := commandValid(executablePath)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("no permission to execute '" + executablePath + "'")
	}
	conn.connectionOK = true
	if eventHandler != nil {
		conn.eventHandler = eventHandler
	}
	go conn.runViaCLI()
	return
}

// SetEventHandler sets a new eventHandler for the connection.
func (conn *Connection) SetEventHandler(eventHandler Event) (err error) {
	conn.eventHandler = eventHandler
	return
}

// SetLogger sets a new logger for the connection.
func (conn *Connection) SetLogger(logger *logrus.Entry) (err error) {
	var loggername string
	loggername = "signalVia"
	switch conn.connectionType {
	case connectionCLI:
		loggername += "CLI"
	case connectionDBus:
		loggername += "DBus"
	}
	conn.logger = logger.WithFields(logrus.Fields{
		"Routine": loggername,
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
