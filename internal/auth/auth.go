package auth

import (
	"context"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis() error {
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
		panic(err)
	}

	fmt.Println(pong)
	return nil
}
