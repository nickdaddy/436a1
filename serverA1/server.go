package main

import (
	"436a1repo/serverA1/serverlib"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var ControllerCreateChan = make(chan string)
var ControllerDeleteChan = make(chan string)
var ControllerListChan = make(chan net.Conn)

type connAndName struct {
}

/*MasterRoomMap is a map of the rooms in the server right now*/

/*ServerInit initializes the server and its threads
TODO: Add Remaning threads needed to be made
*/
func ServerInit() {
	fmt.Printf("Initializing Server\n")
	go Controller()
	ListenerConnThread()
	//fmt.Printf("Server INitialized")

}

/*ListenerConnThread Listens for new connections. Begins a goroutine to handle the new Client connection.
//TODO: Add a channel(?) to know when to get shut down.
https://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/*/
func ListenerConnThread() {
	ln, err := net.Listen("tcp", ":1260")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			os.Exit(1)
		}
		go HandleNewClient(conn)
	}
}

/*HandleNewClient Handles the creation of a new client
TODO: Decide if thread or function; Update to the new structure
*/
func HandleNewClient(conn net.Conn) {
	reader := json.NewDecoder(conn)
	fmt.Println("Recieved new Client")

	for {
		if reader.More() {
			fmt.Println("recieved msg")
			msg := new(serverlib.Message)
			reader.Decode(&msg)
			HandleMessage(msg, conn)
		}
	}

}

/*HandleMessage handles the message*/
func HandleMessage(msg *serverlib.Message, conn net.Conn) {
	//get rid of /n
	msg.Body = msg.Body[0 : len(msg.Body)-1]
	split := strings.SplitN(msg.Body, " ", 3)

	switch split[0] {
	case "/create":
		name := split[1]

		ControllerCreateChan <- name

	case "/delete":
		name := split[1]
		ControllerDeleteChan <- name
	case "/list":
		ControllerListChan <- conn
	case "/join":

	case "/leave":

	}

}

/*Controller adds and removes rooms to the master room map*/
func Controller() {
	MasterRoomMap := make(map[string]*serverlib.Room)

	fmt.Println("Controller ready")
	for {
		select {
		case roomname := <-ControllerCreateChan:
			room := serverlib.NewRoom(roomname)
			MasterRoomMap[room.Name] = room
			go room.StartRoom()
		case roomname := <-ControllerDeleteChan:
			room := MasterRoomMap[roomname]
			room.Deletechan <- true
			delete(MasterRoomMap, roomname)
		case conn := <-ControllerListChan:
			var txt string
			for key := range MasterRoomMap {
				txt = txt + key
			}
			go Send(conn, txt)
			//case roomname := <-ControllerJoinChan:

			//MasterRoomMap[roomname].Joinchan <-
		}
	}

}

/*Send sends back msgs to a specific user, only used for error messages now*/
func Send(conn net.Conn, txt string) {
	msg := new(serverlib.Message)
	msg.Body = txt
	msg.Body = time.Now().String()
	writer := json.NewEncoder(conn)
	writer.Encode(msg)
}

func main() {

	ServerInit()

}
