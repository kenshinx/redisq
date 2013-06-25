package redisq

import (
	"github.com/hoisie/redis"
	"log"
)

type RedisQueue struct {
	redis *redis.Client
	name  string
}

func (rq *RedisQueue) Put(msgs []string) error {
	for _, m := range msgs {
		err := rq.redis.Rpush(rq.name, []byte(m))
		if err != nil {
			log.Printf("insert '%s' into redis queue failed: %s\n", string(m), err)
			return err
		}
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
		log.Printf("get message from redis queue failed: %s\n", err)
	}
	log.Printf("get: %s\n", string(msg))
	return
}

func (rq *RedisQueue) GetUnBlock() (msg []byte, err error) {
	msg, err = rq.Get(false, 0)
	return
}

func (rq *RedisQueue) String() string {
	return "Redis Queue: " + rq.name

}

// Setting password is invalid now. hoisie/redis auth didn't implemented yet.
// Waitting for https://github.com/hoisie/redis/pull/21 merged
func NewRedisQueue(addr string, db int, password string, name string) (rq *RedisQueue) {
	var client = &redis.Client{Addr: addr, Db: db, Password: password}
	rq = &RedisQueue{}
	rq.redis = client
	rq.name = name
	return
}
