package server

import (
	"net"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

type Client struct {
	conn          net.Conn
	inTransaction bool
	queue         []protocol.Value
}

func newTestClient() *Client {
	return &Client{}
}
