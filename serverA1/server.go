package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

/*Message is the message Structure*/
type Message struct {
	Body string
	Date string
}

/*Clients are the current connections being handled
TODO: Change to conn type, or create a new structure with conn in them
*/
var Clients []string

/*ServerInit initializes the server and its threads
TODO: Add Remaning threads needed to be made
*/
func ServerInit() {
	fmt.Printf("Initializing Server\n")
	Clients = make([]string, 0, 10)
	//go ListenerThread()
	go SenderThread()
	ListenerConnThread()

	fmt.Printf("Server INitialized")

}

/*SenderThread blocks on a channel that fills with messages users wish to send across to each other. ListenerMessages thread will handle sent messages and blace them in the channel
TODO: Implement
*/
func SenderThread() {

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
	var found bool

	if len(Clients) == 0 {
		StartNewClient(conn)
		found = true
	} else {
		adrs := conn.RemoteAddr().String()
		for _, element := range Clients {
			if element == adrs {
				found = true
				fmt.Printf("Found Client\n")
				break
			}
		}
	}

	if !found {
		StartNewClient(conn)
	}

	//fmt.Println(msg)

	//defer conn.Close()
	daytime := time.Now().String()
	conn.Write([]byte(daytime)) // don't care about return value
	// we're finished with this client
}

/*StartNewClient is a simple function that adds a new Client called by HandleNewClient*/
func StartNewClient(conn net.Conn) {
	adrs := conn.RemoteAddr().String()
	Clients = append(Clients, adrs)
	var msg = Message{"Welcome to nib851 systems", time.Now().String()}

	e := json.NewEncoder(conn)
	e.Encode(msg)
	fmt.Printf("Client array: %v\n", Clients)
}

func main() {

	ServerInit()

}
