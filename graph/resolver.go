package graph

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/redis/go-redis/v9"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *elasticsearch.Client
	CH *redis.Client
}
