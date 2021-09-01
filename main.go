package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"redlock/db"
)

var (
	instanceNum string
)

func init() {
	instanceNum = os.Getenv("NUM")
}

func main() {
	client := db.NewClient(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}),
		context.Background(),
		instanceNum,
		meow)

	client.ClientFunc()
}

func meow() {
	fmt.Println("мяу", instanceNum)
}
