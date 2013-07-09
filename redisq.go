package redisq

import (
	"encoding/json"
	"github.com/hoisie/redis"
	"log"
	"os"
)

type EmptyQueue struct {
	name string
}

func (e EmptyQueue) Error() string {
	return "Empty Queue: " + e.name

}

type RedisQueue struct {
	redis      *redis.Client
	name       string
	serializer Serializer
	logger     *log.Logger
}

func NewRedisAndQueue(addr string, db int, password string, name string) (rq *RedisQueue) {
	var client = &redis.Client{Addr: addr, Db: db, Password: password}
	rq = &RedisQueue{}
	rq.redis = client
	rq.name = name
	rq.serializer = JsonSerializer{}
	rq.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return
}

func NewRedisQueue(redis *redis.Client, name string) (rq *RedisQueue) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	rq = &RedisQueue{redis, name, JsonSerializer{}, logger}
	return
}

func (rq *RedisQueue) SetLogger(logger *log.Logger) {
	rq.logger = logger
}

func (rq *RedisQueue) Put(msg interface{}) error {

	encoded, err := rq.serializer.Dumps(msg)
	if err != nil {
		rq.logger.Printf("[redis queue] %v encode failed:%s", msg, err)
		return err
	}
	if err := rq.redis.Rpush(rq.name, encoded); err != nil {
		rq.logger.Printf("[redis queue] insert '%v' failed: %s\n", encoded, err)
		return err
	}

	return nil
}

func (rq *RedisQueue) Get(block bool, timeout uint) (i interface{}, err error) {
	var msg []byte
	if block {
		_, msg, err = rq.redis.Blpop([]string{rq.name}, timeout)
	} else {
		msg, err = rq.redis.Lpop(rq.name)
	}
	if err != nil {
		rq.logger.Printf("[redis queue] get message failed: %s\n", err)
		return nil, err
	}
	if msg == nil {
		return nil, EmptyQueue{rq.name}
	}
	i, err = rq.serializer.Loads(msg)
	if err != nil {
		rq.logger.Printf("[redis queue] %s decode failed:%s", msg, err)
	}
	return i, nil
}

func (rq *RedisQueue) GetNoWait() (i interface{}, err error) {
	i, err = rq.Get(false, 0)
	return
}

func (rq *RedisQueue) Consume(block bool, timeout uint, msgs chan interface{}) {
	go func() {
		for {
			i, err := rq.Get(block, timeout)
			if err != nil {
				rq.logger.Printf("[redis queue] get message failed. since of: %s\n", err)
				continue
			}
			msgs <- i
		}
	}()
}

func (rq *RedisQueue) Length() int {
	len, err := rq.redis.Llen(rq.name)
	if err != nil {
		rq.logger.Printf("[redis queue] get length failed: %s\n", err)
		return -1
	}
	return len
}

func (rq *RedisQueue) Empty() bool {
	len, err := rq.redis.Llen(rq.name)
	if err != nil {
		rq.logger.Printf("[redis queue] get length failed: %s\n", err)
		return true
	}
	return len == 0
}

func (rq *RedisQueue) Clear() error {
	_, err := rq.redis.Del(rq.name)
	if err != nil {
		rq.logger.Printf("[redis queue ] clear failed : %s\n", err)
		return err
	}
	return nil
}

func (rq *RedisQueue) String() string {
	return "Redis Queue: " + rq.name

}

type Serializer interface {
	Dumps(v interface{}) ([]byte, error)
	Loads(data []byte) (interface{}, error)
}

type JsonSerializer struct {
}

func (JsonSerializer) Dumps(v interface{}) (encoded []byte, err error) {
	encoded, err = json.Marshal(v)
	return
}

func (JsonSerializer) Loads(data []byte) (decoded interface{}, err error) {
	err = json.Unmarshal(data, &decoded)
	return

}
