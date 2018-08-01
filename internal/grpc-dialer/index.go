package grpc_dialer

import (
	"github.com/docker/docker/api/server/router/session"
	"net"
	"google.golang.org/grpc"
)

func dialGRPC() {
	rpcChannel, err := grpc.Dial(session.Addr().String(), grpc.WithInsecure(), grpc.WithDialer(func(s string, duration time.Duration) (net.Conn, error) {
		return session.Open()
	}))
}
