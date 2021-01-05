package signalgo

func (conn *Connection) runViaCLI() {
	if conn.running {
		conn.logger.Debug("Already running")
		return
	}
	if conn.connectionType != connectionCLI {
		conn.logger.Error("This is not a cli connection")
	}
	conn.running = true
	defer func() { conn.running = false }()
	conn.logger.Debug("Started")
	for {
		select {
		case <-conn.eventListener:
			conn.logger.Debug("Stopping")
			return
		}
	}
}
