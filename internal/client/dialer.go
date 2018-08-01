package client

import (
	"net"
	_ "github.com/gongt/simple-yamux/internal"
	"github.com/hashicorp/yamux"
	"log"
	"github.com/gongt/proxy-gateway/api"
	"context"
	"encoding/binary"
	"fmt"
)

type MultiplexClient struct {
	mapper  map[uint32]net.Addr
	session *yamux.Session
}

func NewMultiplexClient(conn net.Conn) *MultiplexClient {
	session, err := yamux.Client(conn, nil)
	if err != nil {
		log.Fatal("create mux client failed: ", err)
	}

	if err != nil {
		log.Fatal("create rpc channel failed: ", err)
	}

	return &MultiplexClient{
		mapper:  make(map[uint32]net.Addr),
		session: session,
	}
}

func (m *MultiplexClient) Open(from, to net.Addr) uint32 {
	if from.Network() == "tcp" {
		log.Println("connecting with OpenTCP.")
		return m.OpenTCP(from.String(), to)
	} else if from.Network() == "tcp" {
		log.Println("connecting with OpenUnix.")
		return m.OpenUnix(from.String(), to)
	}else{
		log.Println("connecting with OpenUnix.")
		return m.OpenUdp(from.String(), to)
	}
}

func (m *MultiplexClient) OpenTCP(remote string, connect net.Addr) uint32 {
}

func (m *MultiplexClient) OpenUnix(remote string, connect net.Addr) uint32 {
}

func (m *MultiplexClient) OpenUdp(remote string, connect net.Addr) uint32 {
}

func (m *MultiplexClient) EventLoop() {
	go func() {
		<-m.session.CloseChan()
		log.Fatal("Error: connection dropped by server.")
	}()

	for {
		conn, err := m.session.Accept()
		if err != nil {
			log.Fatal("can not accept connection: ", err)
		}
		log.Println("got new connection from server.")

		go m.handle(conn)
	}
}

func (m *MultiplexClient) handle(conn net.Conn) {
	var id uint32

	err := binary.Read(conn, binary.LittleEndian, &id)
	if err != nil {
		duplicateMessage(conn, "read typeId failed:", err)
		conn.Close()
		return
	}

	localConnect, ok := m.mapper[id]
	if !ok {
		duplicateMessage(conn, "invalid typeId:", id)
		conn.Close()
		return
	}

	log.Printf("type id is %d\n", id)

	t, err := net_multiplex.Dial(localConnect)
	if err != nil {
		log.Println("failed to connect local:", err)
		fmt.Fprintln(conn, "target connection to %s failed: %s.\n", localConnect.String(), err.Error())
		conn.Close()
		return
	}
	if t == nil {
		log.Println("not connected, but no error.")
		conn.Close()
		return
	}

	log.Println("local connected, bridge has start.")
	net_multiplex.BridgeConnectionSync(t, "local", conn, "remote")
}
func duplicateMessage(conn net.Conn, s ...interface{}) {
	fmt.Fprintln(conn, s...)
	log.Println(s...)
}
