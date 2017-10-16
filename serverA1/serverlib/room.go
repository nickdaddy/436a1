package serverlib

import (
	//"436a1repo/serverA1/client"

	"encoding/json"
	"fmt"
	"net"
	"time"
)

/*Room Structure*/
type Room struct {
	LastUsedTime  time.Time
	Name          string
	cmap          map[net.Conn]bool
	Joinchan      chan net.Conn
	Deletechan    chan bool
	Leavechan     chan net.Conn
	Broadcastchan chan Message
	writer        json.Encoder
}

/*NewRoom makes a room with a name passed to it
Args: name - what the new room's name should be
returns: a new room*/
func NewRoom(name string) *Room {
	/*if MasterRoomMap[name] != nil {
		//Room already exists, handle error
	}*/
	r := new(Room)
	r.Name = name
	r.cmap = make(map[net.Conn]bool)
	r.Joinchan = make(chan net.Conn)
	r.Deletechan = make(chan bool)
	r.Leavechan = make(chan net.Conn)
	r.Broadcastchan = make(chan Message)
	fmt.Println("New room made")
	return r
}

/*StartRoom is the Room Goroutine startup*/
func (r *Room) StartRoom() {
	for {
		select {
		case conn := <-r.Joinchan:
			r.JoinRoom(conn)
		case <-r.Deletechan:
			r = nil
			return
		case conn := <-r.Leavechan:
			r.LeaveRoom(conn)
		case msg := <-r.Broadcastchan:
			r.Broadcast(msg)
		}
	}
}

/*JoinRoom puts a conn into its Clientmap*/
func (r *Room) JoinRoom(c net.Conn) {
	//check if user is in room
	if _, found := r.cmap[c]; found {
		//User already in room so do nothing
		return
	}
	r.cmap[c] = true

}

/*LeaveRoom makes the client leave the room and removes it from the rooms cmap.*/
func (r *Room) LeaveRoom(c net.Conn) {
	//check if user is not in room
	if _, found := r.cmap[c]; !found {
		//if user is not in room do nothing
		return
	}
	delete(r.cmap, c)
}

/*Broadcast sends messages to every user in the room*/
func (r *Room) Broadcast(msg Message) {
	for conn := range r.cmap {
		writer := json.NewEncoder(conn)
		writer.Encode(msg)

	}
}
