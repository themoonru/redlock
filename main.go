package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

const (
	CheckInterval = time.Second * 1
	MasterTtl     = time.Second * 5
	MasterKeyName = "master"
)

var (
	ctx         = context.Background()
	client      *redis.Client
	instanceNum string
	isMaster    bool

	masterAlreadyElectedError = errors.New("master already elected")
)

func init() {
	instanceNum = os.Getenv("NUM")

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func main() {
	waitCount := 0
	for {
		val, err := client.Get(ctx, MasterKeyName).Result()
		// если ошибка получения значения
		if err != nil {
			// если ошибка связана не с отсутствием ключа
			if err != redis.Nil {
				handleError(err)
				continue
			}

			// если ключ master был пустой
			fmt.Println("master not elected")
			if !isMaster && waitCount < 1 {
				fmt.Println("wait previouse master")
				time.Sleep(2 * CheckInterval)
				waitCount++
				continue
			}
			waitCount = 0
			err2 := setMaster(instanceNum)
			if err2 != nil {
				handleError(err)
			}

			continue
		}

		// если успешно считали ключ master и он не пустой
		isMaster = val == instanceNum
		doWork()

		time.Sleep(CheckInterval)
	}
}

// недостучались до redis
func handleError(err error) {
	if err == masterAlreadyElectedError {
		fmt.Println("master already elected")

		return
	}

	isMaster = false
	fmt.Println("Unknown master status. Self fancing...")

	fmt.Println("Error:", err)
}

func setMaster(val string) error {
	txf := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, MasterKeyName, val, MasterTtl)
			return nil
		})

		return err
	}

	err := client.Watch(ctx, txf, MasterKeyName)
	if err != nil {
		if err == redis.TxFailedErr {
			return masterAlreadyElectedError
		}

		return fmt.Errorf("error on watch: %w", err)
	}

	return nil
}

func doWork() {
	if isMaster {
		fmt.Println("I'm master", instanceNum)
	} else {
		fmt.Println("I'm slave", instanceNum)
	}
}
