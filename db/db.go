package db

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type Database struct {
	client      *redis.Client
	instanceNum string
	isMaster    bool
}

var (
	ctx = context.Background()

	masterAlreadyElectedError = errors.New("master already elected")
)

func (db *Database) NewDatabase(address string) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Database{
		client: client,
	}, nil
}
