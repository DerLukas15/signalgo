package signalgo

import (
	"errors"

	dbus "github.com/godbus/dbus/v5"
)

func (conn *Connection) runViaDBus() {
	if conn.running {
		conn.logger.Debug("Already running")
		return
	}
	if conn.connectionType != connectionDBus {
		conn.logger.Error("This is not a dbus connection")
	}
	conn.running = true
	defer func() { conn.running = false }()
	conn.logger.Debug("Adding dbus signal listener")

	err := conn.connection.AddMatchSignal(
		dbus.WithMatchInterface("org.asamk.Signal"),
		dbus.WithMatchObjectPath("/org/asamk/Signal"),
	)
	if err != nil {
		conn.logger.Error(err)
		return
	}
	conn.logger.Debug("Started listening on dbus")

	events := make(chan *dbus.Message, 10)
	conn.connection.Eavesdrop(events)

	for {
		select {
		case <-conn.eventListener:
			conn.logger.Debug("Stopping listening")
			return
		case message := <-events:
			if message != nil && message.IsValid() == nil {
				variant, ok := message.Headers[dbus.FieldMember]
				if ok {
					v := variant.String()
					if len(v) > 0 && v[0] == '"' {
						v = v[1:]
					}
					if len(v) > 0 && v[len(v)-1] == '"' {
						v = v[:len(v)-1]
					}
					switch v {
					case "MessageReceived":
						timestamp := message.Body[0].(int64)
						sender := message.Body[1].(string)
						text := message.Body[3].(string)
						attachments := message.Body[4].([]string)
						//dbus message serial
						//serial := int64(message.Serial())
						if conn.eventHandler != nil {
							go conn.eventHandler.OnMessage(sender, text, attachments, timestamp)
						}
					case "ReceiptReceived":
						timestamp := message.Body[0].(int64)
						sender := message.Body[1].(string)
						//dbus message serial
						//serial := int64(message.Serial())
						if conn.eventHandler != nil {
							go conn.eventHandler.OnMessageRead(sender, timestamp)
						}
					default:
						conn.logger.Debug("Message Type "+v+" not defined: ", message)
					}
				} else {
					replySerial, ok := message.Headers[dbus.FieldReplySerial]
					if ok {
						conn.logger.Debug("Got reply message for serial ", replySerial.String(), ": ", message.Body)
						conn.lock.RLock()
						if replyChan, ok := conn.replies[replySerial.Value().(uint32)]; ok {
							replyChan <- message
						}
						conn.lock.RUnlock()
					}
				}
			}
		}
	}
}

func (conn *Connection) sendMessageDBus(target, message string, attachments []string) (int64, error) {
	if !conn.running {
		return 0, errors.New("Not running")
	}
	if !conn.connectionOK {
		return 0, errors.New("connection not valid")
	}
	if conn.connectionType != connectionDBus {
		return 0, errors.New("This is not a dbus connection")
	}
	if target == "" || message == "" && len(attachments) == 0 {
		return 0, errors.New("No Target or Message set")
	}
	conn.logger.Debug("Sending Message")
	var body []interface{}
	body = append(body, message)
	body = append(body, attachments)
	body = append(body, target)
	msg := new(dbus.Message)
	msg.Type = dbus.TypeMethodCall
	msg.Flags = 0
	msg.Headers = make(map[dbus.HeaderField]dbus.Variant)
	msg.Headers[dbus.FieldPath] = dbus.MakeVariant(dbus.ObjectPath("/org/asamk/Signal"))
	msg.Headers[dbus.FieldDestination] = dbus.MakeVariant("org.asamk.Signal")
	msg.Headers[dbus.FieldMember] = dbus.MakeVariant("sendMessage")
	msg.Headers[dbus.FieldInterface] = dbus.MakeVariant("org.asamk.Signal")
	msg.Body = body
	if len(body) > 0 {
		msg.Headers[dbus.FieldSignature] = dbus.MakeVariant(dbus.SignatureOf(body...))
	}
	call := conn.connection.Send(msg, nil)
	//Clear Call Thread
	call.ContextCancel()
	if call.Err != nil {
		return 0, call.Err
	}
	serial := msg.Serial()
	conn.logger.Debug("Message DBus serial: ", serial)
	replyChan := make(chan *dbus.Message)
	conn.lock.Lock()
	conn.replies[serial] = replyChan
	conn.lock.Unlock()
	reply := <-replyChan
	messageTimestamp := reply.Body[0].(int64)
	return messageTimestamp, nil
}
