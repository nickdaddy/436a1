package main

import (
	"436a1repo/serverA1/serverlib"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

/*various channels for interacting with controller*/
var controllerCreateChan = make(chan *serverlib.ConnPack)
var controllerDeleteChan = make(chan *serverlib.ConnPack)
var controllerListChan = make(chan net.Conn)
var controllerJoinChan = make(chan *serverlib.ConnPack)
var controllerLeaveChan = make(chan *serverlib.ConnPack)
var controllerBroadcastChan = make(chan *serverlib.ConnPack)
var controllerQuitChan = make(chan net.Conn)
var port string
var numclients = 0

var cap = 10
var numclientsmtx sync.Mutex

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
 */
func ListenerConnThread() {
	ln, err := net.Listen("tcp", port)
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

		numclientsmtx.Lock()
		if numclients < cap {
			go HandleNewClient(conn)
			numclients++
		} else {
			Send(conn, "Server is full please try again later")
			Send(conn, "exit")
			conn.Close()
		}
		numclientsmtx.Unlock()
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
	split := strings.SplitN(msg.Body, " ", 2)
	//fmt.Println("SPlit 0: " + split[0] + " split1: " + split[1])
	switch split[0] {
	case "/create":
		if len(split) > 1 {
			cp := new(serverlib.ConnPack)
			cp.Conn = conn
			cp.Rname = split[1]

			controllerCreateChan <- cp
		} else {

			errormsg := "Room not created: Must invoke with /create roomname"
			go Send(conn, errormsg)
		}

	case "/delete":
		if len(split) > 1 {
			cp := serverlib.NewConnPack(conn, split[1], "")
			controllerDeleteChan <- cp
		} else {
			errormsg := "ServerError: Room not delete: Must invoke with /delete roomname"
			go Send(conn, errormsg)
		}

	case "/list":
		controllerListChan <- conn
	case "/join":
		if len(split) > 1 {
			cp := serverlib.NewConnPack(conn, split[1], "")
			controllerJoinChan <- cp
		} else {
			errormsg := "ServerError: Room not Joined: Must invoke with /join roomname."
			go Send(conn, errormsg)
		}

	case "/leave":
		if len(split) > 1 {
			cp := serverlib.NewConnPack(conn, split[1], "")
			controllerLeaveChan <- cp
		} else {
			errormsg := "ServerError: Room not left: not enough args: Must invoke with /join roomname."
			go Send(conn, errormsg)
		}

	case "/quit":
		controllerQuitChan <- conn

	default:
		if len(split) > 1 {
			crn := new(serverlib.ConnPack)
			crn.Rname = split[0][1:len(split[0])]
			crn.Conn = conn
			msg.Name = msg.Name[0 : len(msg.Name)-1]
			crn.Txt = msg.Name + " : " + split[1]
			controllerBroadcastChan <- crn
		} else {
			err := "ServerError: Message can't be interpreted. You may invoke with /roomname My message here"
			go Send(conn, err)
		}
	}

}

/*Controller handles requests passed from each clients HandleMessage go routine. There is only 1 controller which allows for concurrency*/
func Controller() {
	MasterRoomMap := make(map[string]*serverlib.Room)

	fmt.Println("Controller ready")
	for {
		select {
		case cp := <-controllerCreateChan:
			if _, found := MasterRoomMap[cp.Rname]; found {
				err := "Server Error: Cannot create Room of same name of another room"
				Send(cp.Conn, err)
			} else {
				room := serverlib.NewRoom(cp.Rname)
				MasterRoomMap[room.Name] = room
				go room.StartRoom()
			}
		case cp := <-controllerDeleteChan:
			if _, found := MasterRoomMap[cp.Rname]; found {
				room := MasterRoomMap[cp.Rname]
				room.Deletechan <- true
				delete(MasterRoomMap, cp.Rname)
			} else {
				err := "Server Error: Cannot delete Room that Doesnt exist"
				go Send(cp.Conn, err)
			}

		case cp := <-controllerJoinChan:
			if room, found := MasterRoomMap[cp.Rname]; found {
				room.Joinchan <- cp.Conn
			} else {

				go Send(cp.Conn, "Server Error: Cannot join: "+cp.Rname+" room Doesnt exist")
			}

		case cp := <-controllerLeaveChan:
			if room, found := MasterRoomMap[cp.Rname]; found {
				room.Leavechan <- cp.Conn
			} else {
				err := "Server Error: Cannot Leave: " + cp.Rname + " Room Doesnt exist"
				go Send(cp.Conn, err)
			}

		case cp := <-controllerBroadcastChan:
			if room, found := MasterRoomMap[cp.Rname]; found {
				room.Broadcastchan <- *cp
			} else {
				err := "Server Error: Room  " + cp.Rname + " was not found. Room Doesnt exist"
				go Send(cp.Conn, err)
			}
		case conn := <-controllerListChan:
			txt := "Available Rooms: "
			for key := range MasterRoomMap {
				txt = txt + "\n" + key
			}

			go Send(conn, txt)

		case conn := <-controllerQuitChan:
			for _, r := range MasterRoomMap {
				r.Leavechan <- conn
			}
			Send(conn, "exit")
			conn.Close()
			numclientsmtx.Lock()
			numclients--
			numclientsmtx.Unlock()
		default:

			for rname, r := range MasterRoomMap {
				if time.Since(r.LastUsedTime).Seconds() > (time.Duration(60) * time.Second).Seconds() {
					msg := serverlib.NewMessage("Room: " + rname + "has been inactive for too long")
					r.BroadcastAll(msg)
					delete(MasterRoomMap, rname)
					r.Deletechan <- true

				}
			}

		}
	}

}

/*Send sends back msgs to a specific user, only used for error messages now*/
func Send(conn net.Conn, txt string) {
	writer := json.NewEncoder(conn)

	msg := new(serverlib.Message)
	msg.Body = txt
	msg.Date = time.Now()

	writer.Encode(msg)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Improper arguments, please invoke with ./serverA1 portnumber")
		os.Exit(2)
	}

	if os.Args[1] == "-t" {
		port = ":1260"
	} else {
		port = ":" + os.Args[1]

	}

	ServerInit()

}
