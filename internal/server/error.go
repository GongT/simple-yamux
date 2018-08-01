package server

import (
	"net"
	"log"
)

func handleMultiplexError(err error, conn net.Conn) bool {
	if err != nil {
		log.Println("multiplex error: " + err.Error())
		conn.Close()
		return true
	}
	return false
}
