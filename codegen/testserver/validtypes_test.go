package testserver

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinhtran24/gqlgen/client"
	"github.com/tinhtran24/gqlgen/handler"
)

func TestValidType(t *testing.T) {
	resolvers := &Stub{}
	resolvers.QueryResolver.ValidType = func(ctx context.Context) (validType *ValidType, e error) {
		return &ValidType{
			DifferentCase:    "new",
			DifferentCaseOld: "old",
		}, nil
	}

	srv := httptest.NewServer(handler.GraphQL(NewExecutableSchema(Config{Resolvers: resolvers})))
	c := client.New(srv.URL)

	t.Run("fields with differing cases can be distinguished", func(t *testing.T) {
		var resp struct {
			ValidType struct {
				New string `json:"differentCase"`
				Old string `json:"different_case"`
			}
		}
		err := c.Post(`query { validType { differentCase, different_case } }`, &resp)
		require.NoError(t, err)

		require.Equal(t, "new", resp.ValidType.New)
		require.Equal(t, "old", resp.ValidType.Old)
	})
}
