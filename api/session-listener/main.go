package session_listener

import (
	"net"
	"github.com/pkg/errors"
	"context"
	"github.com/gongt/simple-yamux/api"
	"sync"
	"time"
)

type SessionListener struct {
	l1       net.Listener
	l2       net.Listener
	ctx      context.Context
	address  net.Addr
	closed   bool
	cancel   context.CancelFunc
	sessions map[net.Conn]net.Conn
	mu       sync.Mutex
}

func Listen(network api.NetworkType, address string) (ret *SessionListener, err error) {
	ret = &SessionListener{
		sessions: make(map[net.Conn]net.Conn),
	}

	switch network {
	case api.NetworkTcp:
		ret.address, err = net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return
		}
		ret.l1, err = net.ListenTCP("tcp4", ret.address.(*net.TCPAddr))
		if err != nil {
			return
		}
		ret.l2, err = net.ListenTCP("tcp6", ret.address.(*net.TCPAddr))
		if err != nil {
			ret.l2 = nil
			err = nil
		}
	case api.NetworkUnix:
		ret.address, err = net.ResolveUnixAddr("unix", address)
		if err != nil {
			return
		}
		ret.l1, err = net.ListenUnix("tcp4", ret.address.(*net.UnixAddr))
		if err != nil {
			return
		}
	case api.NetworkKcp:
		err = errors.New("not impl network type: " + network)
	default:
		err = errors.New("unknown network type: " + network + ". available: tcp, unix, kcp")
		return
	}

	return
}

func (sl *SessionListener) Run(ctx context.Context, cbConn api.ConnectionHandler, cbErr api.ErrorHandler) error {
	if sl.ctx != nil {
		return errors.New("session listener is already started")
	}

	sl.ctx, sl.cancel = context.WithCancel(ctx)

	connCh, errCh := makeAcceptChan()
	go popAccept(ctx, sl.l1, connCh, errCh)

	if sl.l2 != nil {
		go popAccept(ctx, sl.l2, connCh, errCh)
	}

	for {
		select {
		case e := <-errCh:
			if cbErr != nil {
				cbErr(e)
			}
		case c := <-connCh:
			sl.register(c)
			cbConn(c)
			sl.unregister(c)
			c.Close()
		case <-ctx.Done():
			break
		}
	}

	sl.closed = true
	sl.ctx = nil

	close(connCh)
	close(errCh)
	return nil
}

func (sl *SessionListener) StopListen() error {
	if sl.ctx != nil {
		return errors.New("session listener is not started")
	}
	if sl.closed {
		return errors.New("session listener is already closed")
	}

	sl.cancel()

	return nil
}

func (sl *SessionListener) Close() error { //graceful
	err := sl.StopListen()
	if err != nil {
		return err
	}
	sl.WaitAll()
	return nil
}

func (sl SessionListener) WaitAll() {
	if sl.ctx != nil {
		<-sl.ctx.Done()
	}

	for {
		if len(sl.sessions) == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (sl *SessionListener) register(conn net.Conn) {
	sl.mu.Lock()
	sl.sessions[conn] = conn
	sl.mu.Unlock()
}
func (sl *SessionListener) unregister(conn net.Conn) {
	sl.mu.Lock()
	delete(sl.sessions, conn)
	sl.mu.Unlock()
}
