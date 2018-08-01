package server

import (
	"net"
	"github.com/hashicorp/yamux"
	"github.com/gongt/proxy-gateway/internal/rpc-server"
	_ "github.com/gongt/simple-yamux/internal"
	"log"
	"encoding/binary"
	"fmt"
	"github.com/gongt/simple-yamux/internal"
)

type MultiplexServer struct {
	conn    net.Conn
	session *yamux.Session
}

func NewServer(conn net.Conn) (ret *MultiplexServer, err error) {
	session, err := yamux.Server(conn, internal.DefaultConfig())
	if err != nil {
		return
	}

	ret = &MultiplexServer{
		conn:    conn,
		session: session,
	}

	return
}

func (m *MultiplexServer) Run() {

}

func (m *MultiplexServer) Close() error {
	e := m.conn.Close()
	if e != nil {
		return e
	}
	e = m.session.Close()
	if e != nil {
		return e
	}
	return nil
}

func (m *MultiplexServer) handleChannel(target rpc_server.ConnectionTarget) {
	conn, err := m.session.Open()
	if err != nil {
		log.Println("open sub connection failed:", err)
		return
	}

	log.Printf("guid: %v", uint32(target.Guid))
	err = binary.Write(conn, binary.LittleEndian, uint32(target.Guid))
	if err != nil {
		fmt.Fprintln(conn, "write id failed:", err)
		conn.Close()
		return
	}
}
