package graph

import (
	"context"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	redis9 "github.com/redis/go-redis/v9"
)

func NewClient() *elasticsearch7.Client { //must be capitlized to be visible by other packages

	goClient, err := elasticsearch7.NewDefaultClient()

	if err != nil {

		panic(err)
	}

	return goClient
}

func NewRedisClient() *redis9.Client {

	redisClient := redis9.NewClient(&redis9.Options{
		Addr:     "localhost:6379",
		Password: "", //no password
		DB:       0,  //use default DB
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil { //if cannot reach db return empty client

		return redis9.NewClient(&redis9.Options{
			Addr:     "",
			Password: "",
			DB:       0,
		})
	}

	return redisClient
}
