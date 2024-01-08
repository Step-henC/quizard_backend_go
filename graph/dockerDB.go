package graph

import (
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	redis9 "github.com/redis/go-redis/v9"
)

func NewClient() *elasticsearch7.Client { //must be capitlized to be visible by other packages

	goClient, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: []string{"http://elasticsearch:9200", "http://elasticsearch:9201"},
		Username:  "",
		Password:  "",
	})

	if err != nil {

		panic(err)
	}

	return goClient
}

func NewRedisClient() *redis9.Client {

	redisClient := redis9.NewClient(&redis9.Options{
		Addr:     "redis:6379",
		Password: "", //no password
		DB:       0,  //use default DB
	})

	return redisClient
}
