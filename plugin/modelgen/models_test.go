package modelgen

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinhtran24/gqlgen/codegen/config"
)

func TestModelGeneration(t *testing.T) {
	cfg, err := config.LoadConfig("testdata/gqlgen.yml")
	require.NoError(t, err)
	p := Plugin{}
	require.NoError(t, p.MutateConfig(cfg))

	require.True(t, cfg.Models.UserDefined("MissingType"))
	require.True(t, cfg.Models.UserDefined("MissingEnum"))
	require.True(t, cfg.Models.UserDefined("MissingUnion"))
	require.True(t, cfg.Models.UserDefined("MissingInterface"))
}
