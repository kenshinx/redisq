package main

import (
	"fmt"
	"redisq"
)

const (
	Server   = "127.0.0.1:6379"
	Db       = 0
	Password = ""
	Name     = "redisq:kenshin"
)

var rq = redisq.NewRedisQueue(Server, Db, Password, Name)

func put() {
	mesg := []string{"a", "b", "c", "d"}
	rq.Put(mesg)
}

func get() {
	msg, _ := rq.Get(true, 1)
	fmt.Printf("get: %s\n", msg)
}

func getNoWait() {
	msg, _ := rq.GetNoWait()
	fmt.Printf("get: %s\n", msg)
}

func consume() {
	var msgs = make(chan []byte)
	rq.Consume(true, 1, msgs)
	for {
		msg := <-msgs
		if msg != nil {
			fmt.Printf("get msg: %s\n", string(msg))
		}

	}
}

func clear() {
	rq.Clear()
}

func length() {
	len := rq.Length()
	fmt.Printf("queue length:%d \n", len)
}

func empty() {
	isEmpty := rq.Empty()
	fmt.Printf("queue is empty? :%t\n", isEmpty)
}

func main() {
	clear()
	put()
}
