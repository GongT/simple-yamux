package server

import (
	"github.com/gongt/simple-yamux/api"
	"github.com/gongt/simple-yamux/api/session-listener"
	"context"
	"log"
	"net"
	"github.com/gongt/simple-yamux/internal/server"
)

type Server struct {
	lis     *session_listener.SessionListener
	network api.NetworkType
	addr    string
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewServer(network api.NetworkType, addr string) (s *Server, err error) {
	lis, err := session_listener.Listen(network, addr)

	s = &Server{
		network: network,
		addr:    addr,
		lis:     lis,
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	return
}

func (s *Server) StartListen() error {
	return s.lis.Run(s.ctx, s.handleConn, func(e error) {
		log.Println("listen error: ", e)
	})
}

func (s *Server) handleConn(conn net.Conn) {
	mux, err := server.NewServer(conn)
	if err != nil {
		conn.Close()
		return
	}

	mux.Run()
}
