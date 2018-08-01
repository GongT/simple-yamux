package internal

import "net"

type ConnectionHandler = func(conn net.Conn)
