package serverlib

import (
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
	Broadcastchan chan ConnPack
	writer        json.Encoder
	log           string
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
	r.Broadcastchan = make(chan ConnPack)
	r.setLastUsed()
	fmt.Println("New room made")
	return r
}

/*StartRoom is the Room Goroutine startup*/
func (r *Room) StartRoom() {
	for {
		select {

		case conn := <-r.Joinchan:
			r.joinRoom(conn)
			r.setLastUsed()
		case <-r.Deletechan:

			msg := NewMessage("You are no longer apart of room: " + r.Name + ". Room is being deleted")
			r.BroadcastAll(msg)

			r = nil
			return

		case conn := <-r.Leavechan:
			r.leaveRoom(conn)
			r.setLastUsed()
		case cp := <-r.Broadcastchan:
			m := NewMessage(cp.Txt)
			r.BroadcastAll(m)
			r.setLastUsed()

		}
	}
}

/*JoinRoom puts a conn into its Clientmap*/
func (r *Room) joinRoom(c net.Conn) {
	//check if user is in room
	if _, found := r.cmap[c]; found {
		//User already in room so do nothing
		r.sendTo(c, "You are already in room : "+r.Name)

		return
	}
	txt := "Welcome to Room: " + r.Name + "\nPrevious messages : \n" + r.log
	r.sendTo(c, txt)
	r.cmap[c] = true

}

/*LeaveRoom makes the client leave the room and removes it from the rooms cmap.*/
func (r *Room) leaveRoom(c net.Conn) {
	//check if user is not in room
	if _, found := r.cmap[c]; !found {
		//if user is not in room do nothing
		r.sendTo(c, "You are not in room: "+r.Name)
		return
	}
	r.sendTo(c, "You have left Room : "+r.Name)

	delete(r.cmap, c)
}

/*Broadcast sends messages to every user in the room except for the one who sent the message*/
func (r *Room) Broadcast(c ConnPack) {
	r.log += "\n" + c.Txt

	for conn := range r.cmap {
		if conn == c.Conn {
			continue
		}
		writer := json.NewEncoder(conn)
		msg := NewMessage(c.Txt)
		writer.Encode(msg)

	}
}

/*senTo sends to specific client conn*/
func (r *Room) sendTo(conn net.Conn, txt string) {
	msg := new(Message)
	msg.Body = txt
	msg.Date = time.Now()
	writer := json.NewEncoder(conn)
	writer.Encode(msg)
}

/*BroadcastAll sends message to everyone in room*/
func (r *Room) BroadcastAll(m Message) {
	r.log += "\n" + m.Body
	m.Body = r.Name + ">>> " + m.Body
	for conn := range r.cmap {

		writer := json.NewEncoder(conn)

		writer.Encode(m)

	}
}

/*setLastUsed sets the rooms last used time*/
func (r *Room) setLastUsed() {
	r.LastUsedTime = time.Now()
}
