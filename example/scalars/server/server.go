package main

import (
	"log"
	"net/http"

	"github.com/tinhtran24/gqlgen/example/scalars"
	"github.com/tinhtran24/gqlgen/handler"
)

func main() {
	http.Handle("/", handler.Playground("Starwars", "/query"))
	http.Handle("/query", handler.GraphQL(scalars.NewExecutableSchema(scalars.Config{Resolvers: &scalars.Resolver{}})))

	log.Fatal(http.ListenAndServe(":8084", nil))
}
