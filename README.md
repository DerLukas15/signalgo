[![Go Reference](https://pkg.go.dev/badge/github.com/DerLukas15/signalgo.svg)](https://pkg.go.dev/github.com/DerLukas15/signalgo)
![GitHub](https://img.shields.io/github/license/DerLukas15/signalgo)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/DerLukas15/signalgo/main)
![GitHub last commit](https://img.shields.io/github/last-commit/DerLukas15/signalgo)
![GitHub top language](https://img.shields.io/github/languages/top/DerLukas15/signalgo)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/DerLukas15/signalgo)

# signalgo
A Go interface for communicating over the signal messenger with the help of [AsamK/signal-cli](https://github.com/AsamK/signal-cli).

There are currently some limitations when using the dbus connection method, as the dbus implementation of signal-cli is still experimental. Those are e.g.:
* Registered as normal messages without any text:
	* Remote delete messages
	* Group changes
* Setting registration pin is not supported
* ...

One should use the CLI connection method if those information and features are required.

# Install

Get this package with
```bash
go get gopkg.in/DerLukas15/signalgo.v0
```

Import this package with

```go
import "github.com/DerLukas15/signalgo"
```

# Usage with DBus

Run signal-cli as a deamon.

Create a struct which satisfies the Event interface:

```go
type eventHandler struct {}

func (eH *eventHandler) OnMessage(source, message string, attachments []string, messagetimestamp int64) error {
	fmt.Println("Got Message ", message, " (", messagetimestamp, ") from ", source, " with attachments ", attachments)
	return nil
}

func (eH *eventHandler) OnMessageRead(source string, messagetimestamp int64) error {
	fmt.Println("Message from ", source, " with timestamp ", messagetimestamp, " has been received")
	return nil
}
```

Now Initialize the package with the eventHandler and everything should work if the signal-cli is installed correctly.
```go
func main() {
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
```

See [GoDoc](https://pkg.go.dev/github.com/DerLukas15/signalgo) for examples.
