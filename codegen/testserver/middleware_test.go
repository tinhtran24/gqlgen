package testserver

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/tinhtran24/gqlgen/client"
	"github.com/tinhtran24/gqlgen/graphql"
	"github.com/tinhtran24/gqlgen/handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	resolvers := &Stub{}
	resolvers.QueryResolver.ErrorBubble = func(ctx context.Context) (i *Error, e error) {
		return &Error{ID: "E1234"}, nil
	}

	resolvers.QueryResolver.User = func(ctx context.Context, id int) (user *User, e error) {
		return &User{ID: 1}, nil
	}

	resolvers.UserResolver.Friends = func(ctx context.Context, obj *User) (users []User, e error) {
		return []User{{ID: 1}}, nil
	}

	srv := httptest.NewServer(
		handler.GraphQL(
			NewExecutableSchema(Config{Resolvers: resolvers}),
			handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
				path, _ := ctx.Value("path").([]int)
				return next(context.WithValue(ctx, "path", append(path, 1)))
			}),
			handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
				path, _ := ctx.Value("path").([]int)
				return next(context.WithValue(ctx, "path", append(path, 2)))
			}),
		))

	c := client.New(srv.URL)

	var resp struct {
		User struct {
			ID      int
			Friends []struct {
				ID int
			}
		}
	}

	called := false
	resolvers.UserResolver.Friends = func(ctx context.Context, obj *User) ([]User, error) {
		assert.Equal(t, []int{1, 2, 1, 2}, ctx.Value("path"))
		called = true
		return []User{}, nil
	}

	err := c.Post(`query { user(id: 1) { id, friends { id } } }`, &resp)

	require.NoError(t, err)
	require.True(t, called)

}
