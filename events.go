package signalMessenger

import (
	dbus "github.com/godbus/dbus/v5"
)

func (conn *Connection) run() {
	if conn.running {
		conn.logger.Debug("Already running")
		return
	}
	conn.running = true
	defer func() { conn.running = false }()
	conn.logger.Debug("Adding dbus signal listener")

	conn.eventListener = make(chan bool)
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
