Redis Queue
=====
A simple redis queue written by go.  Inspired by [python-hotqueue](https://github.com/richardhenry/hotqueue)


## Install

>go get github.com/hoisie/redis

>go get github.com/kenshinx/redisq

## Example

#### Init redis queue object

```
import (
	"github.com/kenshinx/redisq"
)

const (
	Server   = "127.0.1.1:6379"
	Db       = 0
	Password = ""
	Name     = "redisq:kenshin"
)
var rq = redisq.NewRedisQueue(Server, Db, Password, Name)
```

#### Put

```
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

f := Foo{"a", 1}
rq.Put(f)
```
* Map type: map[int]string is unsupport since of json library
* struct is supported imperfect now,will be transform into map, when pop out from redis

query by redis-cli:

	redis 127.0.0.1:6379> lrange redisq:kenshin 0 -1
	1) "\"a\""
	2) "1"
	3) "true"
	4) "1.1"
	5) "[\"a\",\"b\",\"c\"]"
	6) "{\"a\":1,\"b\":2}"
	7) "{\"X\":\"a\",\"Y\":1}"

	//json encoded data

#### Get 

Block get. 
```
msg, _ := rq.Get(true, 1)
//argv1: is_block 
//argv2: timeout, if block=false, timeout is invalid
```

Unblock get

```
msg, _ := rq.Get(false, 0)
```

#### Consume

```
var msgs = make(chan interface{})
rq.Consume(true, 1, msgs)
for {
	msg := <-msgs
	if msg != nil {
		fmt.Printf("get msg: %v,type:%s\n", msg, reflect.TypeOf(msg))
	}

}

```

Much more example checkout from example.go

~~[!!] Note: Auth is invalid now.since of hoisie/redis auth didn't implemented yet~~
