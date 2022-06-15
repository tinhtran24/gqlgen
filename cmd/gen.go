package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/tinhtran24/gqlgen/api"
	"github.com/tinhtran24/gqlgen/codegen/config"
	"github.com/urfave/cli"
)

var genCmd = cli.Command{
	Name:  "generate",
	Usage: "generate a graphql server based on schema",
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "verbose, v", Usage: "show logs"},
		cli.StringFlag{Name: "config, c", Usage: "the config filename"},
	},
	Action: func(ctx *cli.Context) {
		var cfg *config.Config
		var err error
		if configFilename := ctx.String("gen"); configFilename != "" {
			cfg, err = config.LoadConfig(configFilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		} else {
			cfg, err = config.LoadConfigFromDefaultLocations()
			if os.IsNotExist(errors.Cause(err)) {
				cfg = config.DefaultConfig()
			} else if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(2)
			}
		}

		if err = api.Generate(cfg); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(3)
		}
	},
}
