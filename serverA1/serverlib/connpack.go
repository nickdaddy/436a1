package serverlib

import "net"

/*ConnPack useful struct for passing connection and roomname*/
type ConnPack struct {
	Conn  net.Conn
	Rname string
	Txt   string
}

/*NewConnPack Constructor*/
func NewConnPack(conn net.Conn, name string, t string) *ConnPack {
	c := new(ConnPack)
	c.Conn = conn
	c.Rname = name
	c.Txt = t
	return c
}
