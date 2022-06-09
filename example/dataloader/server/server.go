package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tinhtran24/gqlgen/example/dataloader"
	"github.com/tinhtran24/gqlgen/graphql/handler"
	"github.com/tinhtran24/gqlgen/graphql/playground"
)

func main() {
	router := chi.NewRouter()
	router.Use(dataloader.LoaderMiddleware)

	router.Handle("/", playground.Handler("Dataloader", "/query"))
	router.Handle("/query", handler.NewDefaultServer(
		dataloader.NewExecutableSchema(dataloader.Config{Resolvers: &dataloader.Resolver{}}),
	))

	log.Println("connect to http://localhost:8082/ for graphql playground")
	log.Fatal(http.ListenAndServe(":8082", router))
}
