package codegen

import (
	"fmt"
	"go/build"
	"go/parser"
	"go/types"
	"os"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/packages"
)

type Build struct {
	PackageName      string
	Objects          Objects
	Inputs           Objects
	Interfaces       []*Interface
	QueryRoot        *Object
	MutationRoot     *Object
	SubscriptionRoot *Object
	SchemaRaw        map[string]string
	SchemaFilename   SchemaFilenames
	Directives       []*Directive
}

type ModelBuild struct {
	PackageName string
	Models      []Model
	Enums       []Enum
}

type ResolverBuild struct {
	PackageName   string
	ResolverType  string
	Objects       Objects
	ResolverFound bool
}

type ServerBuild struct {
	PackageName         string
	ExecPackageName     string
	ResolverPackageName string
}

// Create a list of models that need to be generated
func (cfg *Config) models() (*ModelBuild, error) {
	namedTypes := cfg.buildNamedTypes()

	pkgs := cfg.getPackage()

	cfg.bindTypes(namedTypes, cfg.Model.Dir(), pkgs)

	models, err := cfg.buildModels(namedTypes, pkgs)
	if err != nil {
		return nil, err
	}
	return &ModelBuild{
		PackageName: cfg.Model.Package,
		Models:      models,
		Enums:       cfg.buildEnums(namedTypes),
	}, nil
}

// bind a schema together with some code to generate a Build
func (cfg *Config) resolver() (*ResolverBuild, error) {
	pkgs := cfg.getPackage()

	destDir := cfg.Resolver.Dir()

	namedTypes := cfg.buildNamedTypes()

	cfg.bindTypes(namedTypes, destDir, pkgs)

	objects, err := cfg.buildObjects(namedTypes, pkgs)
	if err != nil {
		return nil, err
	}

	def, _ := findGoType(pkgs, cfg.Resolver.ImportPath(), cfg.Resolver.Type)
	resolverFound := def != nil

	return &ResolverBuild{
		PackageName:   cfg.Resolver.Package,
		Objects:       objects,
		ResolverType:  cfg.Resolver.Type,
		ResolverFound: resolverFound,
	}, nil
}

func (cfg *Config) server(destDir string) *ServerBuild {
	return &ServerBuild{
		PackageName:         cfg.Resolver.Package,
		ExecPackageName:     cfg.Exec.ImportPath(),
		ResolverPackageName: cfg.Resolver.ImportPath(),
	}
}

// bind a schema together with some code to generate a Build
func (cfg *Config) bind() (*Build, error) {
	namedTypes := cfg.buildNamedTypes()

	pkgs := cfg.getPackage()

	cfg.bindTypes(namedTypes, cfg.Exec.Dir(), pkgs)

	objects, err := cfg.buildObjects(namedTypes, pkgs)
	if err != nil {
		return nil, err
	}

	inputs, err := cfg.buildInputs(namedTypes, pkgs)
	if err != nil {
		return nil, err
	}
	directives, err := cfg.buildDirectives(namedTypes)
	if err != nil {
		return nil, err
	}

	b := &Build{
		PackageName:    cfg.Exec.Package,
		Objects:        objects,
		Interfaces:     cfg.buildInterfaces(namedTypes, pkgs),
		Inputs:         inputs,
		SchemaRaw:      cfg.SchemaStr,
		SchemaFilename: cfg.SchemaFilename,
		Directives:     directives,
	}

	if cfg.schema.Query != nil {
		b.QueryRoot = b.Objects.ByName(cfg.schema.Query.Name)
	} else {
		return b, fmt.Errorf("query entry point missing")
	}

	if cfg.schema.Mutation != nil {
		b.MutationRoot = b.Objects.ByName(cfg.schema.Mutation.Name)
	}

	if cfg.schema.Subscription != nil {
		b.SubscriptionRoot = b.Objects.ByName(cfg.schema.Subscription.Name)
	}
	return b, nil
}

func (cfg *Config) validate() error {
	progLoader := cfg.newLoaderWithErrors()
	_, err := progLoader.Load()
	return err
}

var mode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedModule |
	packages.NeedDeps

func (cfg *Config) newLoaderWithErrors() loader.Config {
	conf := loader.Config{ParserMode: parser.ParseComments}
	pkgs, err := packages.Load(&packages.Config{Mode: mode}, cfg.Models.referencedPackages()...)
	if err != nil {
		return conf
	}
	for _, val := range pkgs {
		conf.Import(val.PkgPath)
	}
	return conf
}

func (cfg *Config) getPackage() []*packages.Package {
	pkgs, err := packages.Load(&packages.Config{Mode: mode}, cfg.Models.referencedPackages()...)
	if err != nil {
		return []*packages.Package{}
	}
	return pkgs
}

func (cfg *Config) newLoaderWithoutErrors() loader.Config {
	conf := cfg.newLoaderWithErrors()
	conf.AllowErrors = true
	conf.TypeChecker = types.Config{
		Error: func(e error) {},
	}
	return conf
}

func resolvePkg(pkgName string) (string, error) {
	cwd, err := os.Executable()

	pkg, err := build.Default.Import(pkgName, cwd, build.FindOnly)
	if err != nil {
		return "", err
	}

	return pkg.ImportPath, nil
}
