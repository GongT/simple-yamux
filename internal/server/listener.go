package server

import (
	"net"
	_ "github.com/gongt/simple-yamux/internal"
	"log"
	"encoding/binary"
	"fmt"
	"github.com/gongt/simple-yamux/internal"
)

type MultiplexListener struct {
	lis          net.Listener
	handler      internal.ConnectionHandler
	sessionCount int
}

func NewLis(lis net.Listener, handler internal.ConnectionHandler) (*MultiplexListener) {
	return &MultiplexListener{
		lis:     lis,
		handler: handler,
	}
}

func (m *MultiplexListener) Run() {
	go func() {
		accept := m.WaitRequest()
		for conn := range accept {
			m.handleChannel(conn)
		}
	}()

	prof.Snapshot("new conn")
	server.Start()
	prof.Snapshot("conn fin")
}

func HandleChannel() {
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

	log.Println("open sub connection success. start bridging.")
	net_multiplex.BridgeConnectionSync(target.Conn, "accepted", conn, "remote")
}
