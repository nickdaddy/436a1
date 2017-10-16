package main

import (
	"436a1repo/serverA1/serverlib"
	"bufio"
	"fmt"
	"net"
	"os"
	//"io/ioutil"
	//"time"
	"encoding/json"
)

var Sendchan = make(chan string)

func main() {

	conn, err := net.Dial("tcp", ":1260")
	if err != nil {
		// handle error
		checkError(err)
		os.Exit(1)
	}

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

func ReadStdIn() {
	reader := bufio.NewReader(os.Stdin)

	for {
		txt, _ := reader.ReadString('\n')
		//if err == nil {
		Sendchan <- txt

		//}
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
		}
	}
}

/*Send sends to server*/
func Send(conn net.Conn) {
	encoder := json.NewEncoder(conn)
	for {
		txt := <-Sendchan

		msg := serverlib.NewMessage(txt)

		encoder.Encode(msg)
	}
}
