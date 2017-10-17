package serverlib

import "time"

/*Message struct*/
type Message struct {
	Body string
	Date time.Time
	Name string
}

/*NewMessage makes a new message with the Body of body of type string and the time now*/
func NewMessage(body string) Message {
	m := new(Message)
	m.Body = body
	m.Date = time.Now()
	return *m
}
