// plugin package interfaces are EXPERIMENTAL.

package plugin

import (
	"github.com/tinhtran24/gqlgen/codegen"
	"github.com/tinhtran24/gqlgen/codegen/config"
)

type Plugin interface {
	Name() string
}

type ConfigMutator interface {
	MutateConfig(cfg *config.Config) error
}

type CodeGenerator interface {
	GenerateCode(cfg *codegen.Data) error
}
