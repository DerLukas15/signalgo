# signalMessenger
A Go implementation for communication over the signal messenger via dbus and [AsamK/signal-cli](https://github.com/AsamK/signal-cli)

# Install

Get this package with
```bash
go get gopkg.in/DerLukas15/signalMessenger.v0
```

Import this package with

```go
import "gopkg.in/DerLukas15/signalMessenger.v0"
```

# Usage

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
```

See [Examples](https://github.com/DerLukas15/signalMessenger/_examples) for working implementations.
