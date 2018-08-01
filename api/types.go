package api

import (
	"net"
)

const NetworkTcp = "tcp"
const NetworkUnix = "unix"
const NetworkKcp = "kcp"

type NetworkType = string

type ErrorHandler func(error)
type ConnectionHandler func(conn net.Conn)
