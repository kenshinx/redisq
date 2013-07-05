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

	type Foo struct {
		X string
		Y int
	}

	rq.Put("a")
	rq.Put(1)
	rq.Put(true)
	rq.Put(1.1)
	rq.Put([]string{"a", "b", "c"})
	rq.Put(map[string]int{"a": 1, "b": 2})
	// json unsupported type: map[int]string

	f := Foo{"a", 1}
	rq.Put(f)

	// struct is supported imperfect.
	// strcut type will be transform into map, when pop out from redis.

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
		if msg == nil {
			continue
		}
		fmt.Println(reflect.TypeOf(msg))
		k := reflect.TypeOf(msg).Kind()
		if k >= reflect.Int && k <= reflect.Uint64 {
			fmt.Printf("get data msg:%d\n", msg)
		}
		switch k {
		case reflect.Float32, reflect.Float64:
			fmt.Printf("get float msg:%f\n", msg)
		case reflect.String:
			fmt.Printf("get string msg:%s\n", msg)
		case reflect.Bool:
			fmt.Printf("get bool msg:%t\n", msg)
		case reflect.Slice:
			slice, _ := msg.([]interface{})
			for i, v := range slice {
				fmt.Printf("slice[%d]: %v ,type of %s\n", i, v, reflect.TypeOf(v))
			}
		case reflect.Map:
			tmap, ok := msg.(map[string]interface{})
			if !ok {
				fmt.Printf("bad map struct:%v", tmap)
			}
			for r, v := range tmap {
				fmt.Printf("map[%s]:%v\n", r, v)
			}
		case reflect.Struct:
			//todo
			continue
		default:
			fmt.Printf("Unkown msg:%v\n", msg)
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
