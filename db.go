package main

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type db struct {
	Ctx context.Context
	Client *redis.Client

}

func (d *db) dbInit(){
	d.Ctx = context.Background()
	d.Client = redis.NewClient(&redis.Options{Addr: "localhost:6379", // use default Addr
		Password: "", // no password set
		DB:       0,
	}) // use default DB })
}

func (d *db) dbClose(){
	_ = d.Client.Close()
}
