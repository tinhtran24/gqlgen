package main

import (
	"log"
	"net/http"

	todo "github.com/jlightning/gqlgen/example/config"
	"github.com/jlightning/gqlgen/handler"
)

func main() {
	http.Handle("/", handler.Playground("Todo", "/query"))
	http.Handle("/query", handler.GraphQL(
		todo.NewExecutableSchema(todo.New()),
	))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
