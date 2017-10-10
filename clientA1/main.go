package main

import (
  "net"
  "fmt"
  "os"
  //"io/ioutil"
  //"time"
  "encoding/json"
)

type Message struct {
  Body string
  Date string
}

func main(){

    /*
  if len(os.Args) != 2 {
      fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
      os.Exit(1)
  }

  service := os.Args[1]
  tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
  checkError(err)
  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  checkError(err)
  _, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
  checkError(err)
  result, err := ioutil.ReadAll(conn)
  checkError(err)
  fmt.Println(string(result))
  os.Exit(0)
  */
  conn, err := net.Dial("tcp", ":1260")
  if err != nil {
	// handle error
    checkError(err)
    os.Exit(1)
  }
  //fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
  //status, err := bufio.NewReader(conn).ReadString('\n')

  msg:= Message{
    Body: "",
    Date: "",
  }
  d := json.NewDecoder(conn)

  err = d.Decode(&msg)
  if (err!=nil){
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
  fmt.Println(msg)

}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}



func ListenForMessages(){

}

func SendMessages(){

}
