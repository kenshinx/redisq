package redisq

import (
	"encoding/json"
	"github.com/hoisie/redis"
	"log"
)

type RedisQueue struct {
	redis      *redis.Client
	name       string
	serializer Serializer
}

func NewRedisQueue(addr string, db int, password string, name string) (rq *RedisQueue) {
	var client = &redis.Client{Addr: addr, Db: db, Password: password}
	rq = &RedisQueue{}
	rq.redis = client
	rq.name = name
	rq.serializer = JsonSerializer{}
	return
}

func (rq *RedisQueue) Put(msg interface{}) error {

	encoded, err := rq.serializer.Dumps(msg)
	if err != nil {
		log.Printf("[redis queue] %v encode failed:%s", msg, err)
		return err
	}
	if err := rq.redis.Rpush(rq.name, encoded); err != nil {
		log.Printf("[redis queue] insert '%v' failed: %s\n", encoded, err)
		return err
	}

	return nil
}

func (rq *RedisQueue) Get(block bool, timeout uint) (msg []byte, err error) {
	if block {
		_, msg, err = rq.redis.Blpop([]string{rq.name}, timeout)
	} else {
		msg, err = rq.redis.Lpop(rq.name)
	}
	if err != nil {
		log.Printf("[redis queue] get message failed: %s\n", err)
	}
	return
}

func (rq *RedisQueue) GetNoWait() (msg []byte, err error) {
	msg, err = rq.Get(false, 0)
	return
}

func (rq *RedisQueue) Consume(block bool, timeout uint, msgs chan []byte) {
	go func() {
		for {
			msg, err := rq.Get(block, timeout)
			if err != nil {
				log.Printf("[redis queue] consumer exit. since of: %s\n", err)
				continue
			}
			if msg == nil {
				continue
			}
			msgs <- msg
		}
	}()
}

func (rq *RedisQueue) Length() int {
	len, err := rq.redis.Llen(rq.name)
	if err != nil {
		log.Printf("[redis queue] get length failed: %s\n", err)
		return -1
	}
	return len
}

func (rq *RedisQueue) Empty() bool {
	len, err := rq.redis.Llen(rq.name)
	if err != nil {
		log.Printf("[redis queue] get length failed: %s\n", err)
		return true
	}
	return len == 0
}

func (rq *RedisQueue) Clear() error {
	_, err := rq.redis.Del(rq.name)
	if err != nil {
		log.Printf("[redis queue ] clear failed : %s\n", err)
		return err
	}
	return nil
}

func (rq *RedisQueue) String() string {
	return "Redis Queue: " + rq.name

}

type Serializer interface {
	Dumps(v interface{}) ([]byte, error)
	Loads()
}

type JsonSerializer struct {
}

func (JsonSerializer) Dumps(v interface{}) (encoded []byte, err error) {
	encoded, err = json.Marshal(v)
	return
}

func (JsonSerializer) Loads() {
	return
}
