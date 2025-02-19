package redis

import (
	"github.com/go-redis/redis/v9"
)

type Connection struct {
	Client *redis.Client
}

func Create(Addr, Pass string, idx int) *Connection {
	clinet = redis.NewClient(&redis.Options{
		Addr: 	Addr,
		Password: Pass,
        DB:       idx, 
    })

	return &Connection{Client: client}
}

