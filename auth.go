package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func connect_redis() error {
	pass := os.Getenv("REDISPASS")
	if pass == "" {
		return errors.New("missing redis password")
	}

	client := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		ClientName: "",
		Password:   pass,
		DB:         0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return errors.New("Couldn't connect")
	}

	fmt.Println(pong)
	return nil
}
