package stubgen

import (
	"syscall"

	"github.com/tinhtran24/gqlgen/codegen"
	"github.com/tinhtran24/gqlgen/codegen/config"
	"github.com/tinhtran24/gqlgen/codegen/templates"
	"github.com/tinhtran24/gqlgen/plugin"
)

func New(filename string, typename string) plugin.Plugin {
	return &Plugin{filename: filename, typeName: typename}
}

type Plugin struct {
	filename string
	typeName string
}

var _ plugin.CodeGenerator = &Plugin{}
var _ plugin.ConfigMutator = &Plugin{}

func (m *Plugin) Name() string {
	return "stubgen"
}

func (m *Plugin) MutateConfig(cfg *config.Config) error {
	_ = syscall.Unlink(m.filename)
	return nil
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	return templates.Render(templates.Options{
		PackageName: data.Config.Exec.Package,
		Filename:    m.filename,
		Data: &ResolverBuild{
			Data:     data,
			TypeName: m.typeName,
		},
		GeneratedHeader: true,
	})
}

type ResolverBuild struct {
	*codegen.Data

	TypeName string
}
