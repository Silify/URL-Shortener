package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()

	Client0 *redis.Client
	Client1 *redis.Client
)

func Init() {
	Client0 = newClient(0)
	Client1 = newClient(1)

	if err := Client0.Ping(Ctx).Err(); err != nil {
		log.Fatal("failed to connect to redis db0:", err)
	}
	if err := Client1.Ping(Ctx).Err(); err != nil {
		log.Fatal("failed to connect to redis db1:", err)
	}
}

func newClient(db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       db,
	})
}

func Close() {
	defer Client0.Close()
	defer Client1.Close()
}
