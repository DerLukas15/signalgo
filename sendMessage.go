package signalgo

import (
	"errors"

	dbus "github.com/godbus/dbus/v5"
)

// SendMessage sends a new message to the target with possible attachments.
// This function waits for the message timestamp.
// The timestamp of the message given by signal-cli will be returned.
func (conn *Connection) SendMessage(target, message string, attachments []string) (int64, error) {
	if target == "" || message == "" {
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

// SendMessageASync sends a new message to the target but does not care about the returning messageTimestamp.
func (conn *Connection) SendMessageASync(target, message string, attachments []string) error {
	go conn.SendMessage(target, message, attachments)
	return nil
}
