package extension

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinhtran24/gqlgen/graphql"
)

func TestIntrospection(t *testing.T) {
	rc := &graphql.OperationContext{
		DisableIntrospection: true,
	}
	require.Nil(t, Introspection{}.MutateOperationContext(context.Background(), rc))
	require.Equal(t, false, rc.DisableIntrospection)
}
