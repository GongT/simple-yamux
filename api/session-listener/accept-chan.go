package session_listener

import (
	"net"
	"context"
)

type connChan = chan net.Conn
type errorChan = chan error

func makeAcceptChan() (connChan, errorChan) {
	return make(connChan), make(errorChan)
}

func popAccept(ctx context.Context, lis net.Listener, connChan connChan, errChan errorChan) {
	done := false
	go func() {
		for {
			conn, err := lis.Accept()
			if done {
				return
			}
			if err == nil {
				connChan <- conn
			} else {
				errChan <- err
			}
		}
	}()

	<-ctx.Done()

	done = true
	lis.Close()
}
