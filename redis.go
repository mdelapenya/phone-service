package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var serverURL = os.Getenv("SERVER_URL")

// redisClient Creates a new RedisClient
func redisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     serverURL + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}
