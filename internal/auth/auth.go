package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type data_instance struct {
	cache redis.Options
	db *sqlx.DB
}

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
		panic(err)
	}

	fmt.Println(pong)
	return nil
}

func connect_postgres() *sqlx.DB {
	var sb strings.Builder

	pass := os.Getenv("GOFSDBPASS")

	sb.WriteString("postgres://postgres:")
	sb.WriteString(pass)
	sb.WriteString("@localhost:5432/gofs-db?sslmode=disable")

	db, err := sqlx.Connect("postgres", sb.String())
	if err != nil {
		panic(err)
	}

	return db
	
}
