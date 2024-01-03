package main

import (
	"context"
	"log"

	"os"
	"quizard/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

const defaultPort = "8080"

func GraphqlHandler() gin.HandlerFunc {

	goClient := graph.NewClient()
	dbName := os.Getenv("ELASTIC_DOCKER_DB")
	goClient.Indices.Create(dbName) //create temp database on start up in docker
	goClient.Indices.Create("quizzes")

	redisClient := graph.NewRedisClient()

	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: goClient, CH: redisClient}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := gin.Default()
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, //need this for jwt
		AllowCredentials: true,
		//OptionsPassthrough: true,
		Debug: true,
	})) //line order matters!! Call cors before endpoint set
	r.Use(GinContextToContextMiddleware())

	r.POST("/query", GraphqlHandler())
	r.GET("/", PlaygroundHandler())

	r.Run()

	// srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: goClient}})) //add db here

	//  http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//  http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

}
