package main

import (
	"fmt"
	"redisq"
	"reflect"
)

const (
	Server   = "localhost:6379"
	Db       = 0
	Password = ""
	Name     = "redisq:kenshin"
)

var rq = redisq.NewRedisQueue(Server, Db, Password, Name)

func put() {

	rq.Put("a")
	rq.Put(1)
	rq.Put([]string{"a", "b", "c"})
	rq.Put(map[string]int{"a": 1, "b": 2})
	//json unsupported type: map[int]string

}

func get() {
	msg, _ := rq.Get(true, 1)
	fmt.Printf("get msg: %v,type:%s\n", msg, reflect.TypeOf(msg))
}

func getNoWait() {
	msg, _ := rq.GetNoWait()
	fmt.Printf("get: %v\n", msg)
}

func consume() {
	var msgs = make(chan interface{})
	rq.Consume(true, 1, msgs)
	for {
		msg := <-msgs
		if msg != nil {
			fmt.Printf("get msg: %v,type:%s\n", msg, reflect.TypeOf(msg))
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
	put()
	consume()
}
