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
mesg := []string{"a", "b", "c", "d"}
rq.Put(mesg)
```
*support insert bulk of messages into queue one time*

query by redis-cli:

	redis 127.0.0.1:6379> lrange redisq:kenshin 0 -1
	1) "a"
	2) "b"
	3) "c"
	4) "d"

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
var msgs = make(chan []byte)
rq.Consume(true, 1, msgs)
for {
	msg := <-msgs
	if msg != nil {
		fmt.Printf("get msg: %s\n", string(msg))
	}
}
```

Much more example checkout from example.go

[!!] Note: Auth is invalid now.since of hoisie/redis auth didn't implemented yet

## TODO

* Json encode
* Dynamic type