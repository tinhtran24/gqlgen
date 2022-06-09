package main

import (
	"log"
	"net/http"

	todo "github.com/tinhtran24/gqlgen/example/config"
	"github.com/tinhtran24/gqlgen/graphql/handler"
	"github.com/tinhtran24/gqlgen/graphql/playground"
)

func main() {
	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", handler.NewDefaultServer(
		todo.NewExecutableSchema(todo.New()),
	))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
