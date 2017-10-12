package serverlib

import (
	//"436a1repo/serverA1/client"
	"net"
	"sync"
	"time"
)

/*Room Structure*/
type Room struct {
	LastUsedTime time.Time
	name         string
	cmap         map[net.Conn]*Client
	wait         chan bool
	mux          *sync.Mutex
}

/*NewRoom makes a room with a name passed to it
Args: name - what the new room's name should be
returns: a new room*/
func NewRoom(name string) *Room {
	r := Room{time.Now(), name, make(map[net.Conn]*Client), make(chan bool, 1), &sync.Mutex{}}

	return &r
}

/*EnterThisRoom puts a Client into the room calling it
Args: the client wishing to join the room.
returns: 0 if succesfull, -1 if not. */
func (r Room) EnterThisRoom(c *Client) {

	//	r.wait <- true

	r.mux.Lock()
	c.mux.Lock()
	r.cmap[c.Conn] = c
	c.rmap[r.name] = r
	c.mux.Unlock()
	r.mux.Unlock()
	//<-r.wait

}

/*LeaveRoom makes the client leave the room and dissasociates both maps between the two.*/
func (r Room) LeaveRoom(c *Client) {
	delete(r.cmap, c.Conn)
	//delete(c.rmap, r.Name)

}
