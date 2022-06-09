//go:generate go run ../../../testdata/gqlgen.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tinhtran24/gqlgen/example/federation/products/graph"
	"github.com/tinhtran24/gqlgen/example/federation/products/graph/generated"
	"github.com/tinhtran24/gqlgen/graphql/handler"
	"github.com/tinhtran24/gqlgen/graphql/handler/debug"
	"github.com/tinhtran24/gqlgen/graphql/playground"
)

const defaultPort = "4002"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.Use(&debug.Tracer{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
