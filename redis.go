package main

import (
	"os"

	"github.com/go-redis/redis"
)

var serverURL = os.Getenv("SERVER_URL")

// redisClient Creates a new RedisClient
func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     serverURL + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// Get gets the value of a key from the datastore
func Get(key string) (string, error) {
	client := redisClient()

	return client.Get(key).Result()
}

// Set sets a key and its value to the datastore
func Set(key string, value string) error {
	client := redisClient()

	return client.Set(key, value, 0).Err()
}
