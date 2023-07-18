package task

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}
)

type queue struct {
	name string
}

func (q *queue) Add(name string, args map[string]any) *work.Job {
	enqueuer := work.NewEnqueuer(q.name, redisPool)
	job, err := enqueuer.Enqueue(name, args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return job
}
