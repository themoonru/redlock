package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"redlock/pkg/lock"
)

var (
	instanceNum string
)

func init() {
	instanceNum = os.Getenv("NUM")
}

func main() {
	task := lock.New(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}),
		context.Background(),
		instanceNum,
		meow)

	task.Run()
}

func meow() {
	fmt.Println("мяу", instanceNum)
}
