// Copyright 2021 Lukas Gallandi. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package signalgo_test

import (
	"errors"

	"github.com/DerLukas15/signalgo"
)

func ExampleNewDBus_systemBus() {
	conn, err := signalgo.NewDBus(true, nil)
	if err != nil {
		panic(err)
	}
	if conn == nil {
		panic(errors.New("No connection"))
	}
}

func ExampleNewDBus_sessionBus() {
	conn, err := signalgo.NewDBus(false, nil)
	if err != nil {
		panic(err)
	}
	if conn == nil {
		panic(errors.New("No connection"))
	}
}

func ExampleNewCLI() {
	conn, err := signalgo.NewCLI("PathToExecutable", nil)
	if err != nil {
		panic(err)
	}
	if conn == nil {
		panic(errors.New("No connection"))
	}
}
