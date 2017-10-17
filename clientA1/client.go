package main

import (
	"436a1repo/serverA1/serverlib"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

var sendchan = make(chan string)
var name string

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Improper arguments, please invoke with ./clientA1 hostip:portnumber or -t for testing")
		os.Exit(2)
	}
	address := os.Args[1]
	if os.Args[1] == "-t" {
		address = ":1260"
	}

	fmt.Println("Input a name for yourself")
	reader := bufio.NewReader(os.Stdin)
	name, _ = reader.ReadString('\n')

	conn, err := net.Dial("tcp", address)
	if err != nil {
		// handle error
		checkError(err)
		os.Exit(1)
	}

	fmt.Println("Hello welcome to 436a1NickGoChat with golang")
	fmt.Println(`please use 436a1NickGoChat with the following commands:
/create			-->	creates a room
/join				-->	joins a room
/leave 			-->	leaves a room
/delete			--> deletes a room
/list				--> lists available rooms
/<roomname> --> where roomname is a room in the server. Do not include <>
/quit				--> quit this program and leave all chat rooms
/help				--> displays this message again`)

	go ReadStdIn()
	go Listen(conn)
	Send(conn)

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

/*ReadStdIn reads what the user types into stdin*/
func ReadStdIn() {
	reader := bufio.NewReader(os.Stdin)

	for {
		txt, _ := reader.ReadString('\n')

		if txt[0] != '/' {
			fmt.Println("Improper use of 436a1 Chat: all commands must be invoked with a '/' at the start. input /help if you need to be reminded of commands")
		} else if txt == "/help\n" {
			fmt.Println(`please use 436a1NickGoChat with the following commands:
/create			-->	creates a room
/join 			-->	joins a room
/leave 			-->	leaves a room
/delete 		--> deletes a room
/list				--> lists available rooms
/<roomname> --> where roomname is a room in the server. Do not include <>
/quit				--> quit this program and leave all chat rooms
/help				--> displays this message again`)
		} else if txt == "/quit\n" {
			fmt.Println("Shutting down ClientA1")
			sendchan <- txt
		} else {
			sendchan <- txt

		}

	}
}

/*Listen listens from server*/
func Listen(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	for {
		if decoder.More() {
			msg := new(serverlib.Message)
			decoder.Decode(msg)
			fmt.Println(msg.Body)
			if msg.Body == "exit" {
				os.Exit(0)
			}
		}
	}
}

/*Send sends to server*/
func Send(conn net.Conn) {
	encoder := json.NewEncoder(conn)
	for {
		txt := <-sendchan

		msg := serverlib.NewMessage(txt)
		msg.Name = name
		encoder.Encode(msg)
	}
}
