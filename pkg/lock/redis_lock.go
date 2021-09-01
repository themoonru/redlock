package lock

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	CheckInterval = time.Second * 1
	MasterTtl     = time.Second * 5
	MasterKeyName = "master"
)

var (
	masterAlreadyElectedError = errors.New("master already elected")
)

type RedisDistributedLock struct {
	client      *redis.Client
	ctx         context.Context
	f           func()
	isMaster    bool
	instanceNum string
}

func New(client *redis.Client, ctx context.Context, instanceNum string, f func()) *RedisDistributedLock {
	return &RedisDistributedLock{
		client:      client,
		ctx:         ctx,
		instanceNum: instanceNum,
		f:           f,
	}
}

func (c *RedisDistributedLock) handleError(err error) {
	if err == masterAlreadyElectedError {
		fmt.Println("master already elected")

		return
	}

	c.isMaster = false
	fmt.Println("Unknown master status. Fancing myself...")
	time.Sleep(time.Second)

	fmt.Println("Error:", err)
}

func (c *RedisDistributedLock) setMaster(val string) error {
	txf := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(c.ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(c.ctx, MasterKeyName, val, MasterTtl)
			return nil
		})

		return err
	}

	err := c.client.Watch(c.ctx, txf, MasterKeyName)
	if err != nil {
		if err == redis.TxFailedErr {
			return masterAlreadyElectedError
		}

		return fmt.Errorf("error on watch: %w", err)
	}

	return nil
}

func (c *RedisDistributedLock) Run() {
	waitCount := 0
	for {
		val, err := c.client.Get(c.ctx, MasterKeyName).Result()
		// если ошибка получения значения
		if err != nil {
			// если ошибка связана не с отсутствием ключа
			if err != redis.Nil {
				c.handleError(err)
				continue
			}

			// если ключ master был пустой
			fmt.Println("master not elected")
			if !c.isMaster && waitCount < 1 {
				fmt.Println("wait previous master")
				time.Sleep(2 * CheckInterval)
				waitCount++
				continue
			}
			waitCount = 0
			err2 := c.setMaster(c.instanceNum)
			if err2 != nil {
				c.handleError(err)
			}

			continue
		}

		// если успешно считали ключ master и он не пустой
		c.isMaster = val == c.instanceNum
		c.doMasterWork()

		time.Sleep(CheckInterval)
	}
}

func (c *RedisDistributedLock) doMasterWork() {
	if c.isMaster {
		c.f()
	} else {
		fmt.Println("I'm slave")
	}
}
