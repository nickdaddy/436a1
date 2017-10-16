package main

import (
	"436a1repo/serverA1/serverlib"
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
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
		fmt.Println("client read " + txt)
		Sendchan <- txt

		//}
	}
}

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

func Send(conn net.Conn) {
	encoder := json.NewEncoder(conn)
	for {
		txt := <-Sendchan
		fmt.Println("Sent msg to server")

		msg := new(serverlib.Message)
		msg.Body = txt
		msg.Date = time.Now().String()
		encoder.Encode(msg)
	}
}
