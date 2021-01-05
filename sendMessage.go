package signalgo

import "errors"

// SendMessage sends a new message to the target with possible attachments.
// This function waits for the message timestamp.
// The timestamp of the message given by signal-cli will be returned.
func (conn *Connection) SendMessage(target, message string, attachments []string) (int64, error) {
	switch conn.connectionType {
	case connectionDBus:
		return conn.sendMessageDBus(target, message, attachments)
	default:
		return 0, errors.New("connection type not defined")
	}
}

// SendMessageASync sends a new message to the target but does not care about the returning messageTimestamp.
func (conn *Connection) SendMessageASync(target, message string, attachments []string) error {
	go conn.SendMessage(target, message, attachments)
	return nil
}
