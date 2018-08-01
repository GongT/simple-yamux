package internal

import (
	"github.com/hashicorp/yamux"
	"os"
	"time"
)

func SetKeepAlive(keepalive time.Duration) {
	cfg.KeepAliveInterval = keepalive
	cfg.ConnectionWriteTimeout = keepalive / 2
}

var cfg *yamux.Config

func init() {
	cfg = yamux.DefaultConfig()
	cfg.EnableKeepAlive = true
	SetKeepAlive(10 * time.Second)
	cfg.LogOutput = os.Stderr
}

func DefaultConfig() *yamux.Config {
	return cfg
}
