package serverlib

import (
	"net"
	"sync"
)

/*Client structure*/
type Client struct {
	Name string
	Conn net.Conn
	rmap map[string]Room
	mux  sync.Mutex
}

/*NewClient creates a new client
return: a new client*/
func NewClient(conn net.Conn) *Client {
	c := Client{"fillername", conn, make(map[string]Room), sync.Mutex{}}
	return &c
}
