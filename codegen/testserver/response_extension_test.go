package testserver

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinhtran24/gqlgen/client"
	"github.com/tinhtran24/gqlgen/graphql"
	"github.com/tinhtran24/gqlgen/handler"
)

func TestResponseExtension(t *testing.T) {
	resolvers := &Stub{}
	resolvers.QueryResolver.Valid = func(ctx context.Context) (s string, e error) {
		return "Ok", nil
	}

	srv := httptest.NewServer(handler.GraphQL(
		NewExecutableSchema(Config{Resolvers: resolvers}),
		handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
			rctx := graphql.GetRequestContext(ctx)
			if err := rctx.RegisterExtension("example", "value"); err != nil {
				panic(err)
			}
			return next(ctx)
		}),
	))
	c := client.New(srv.URL)

	raw, _ := c.RawPost(`query { valid }`)
	require.Equal(t, raw.Extensions["example"], "value")
}
