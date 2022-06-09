package testserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinhtran24/gqlgen/client"
	"github.com/tinhtran24/gqlgen/graphql"
	"github.com/tinhtran24/gqlgen/graphql/handler"
)

func TestResponseExtension(t *testing.T) {
	resolvers := &Stub{}
	resolvers.QueryResolver.Valid = func(ctx context.Context) (s string, e error) {
		return "Ok", nil
	}

	srv := handler.NewDefaultServer(
		NewExecutableSchema(Config{Resolvers: resolvers}),
	)

	srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		graphql.RegisterExtension(ctx, "example", "value")

		return next(ctx)
	})

	c := client.New(srv)

	raw, _ := c.RawPost(`query { valid }`)
	require.Equal(t, raw.Extensions["example"], "value")
}
