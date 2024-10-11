package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Piszmog/gtpm/log"
	"github.com/Piszmog/gtpm/run"
	"github.com/urfave/cli/v2"
)

func main() {
	var logger *slog.Logger
	app := &cli.App{
		Usage: "TPM Plugin Manager",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "level",
				Value: "info",
				Usage: "Change the log level (e.g. debug, warn, info)",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "clean",
				Aliases: []string{"c"},
				Usage:   "Clean Plugins",
				Action: func(ctx *cli.Context) error {
					logger = log.New(log.Level(ctx.String("level")), log.OutputText)
					return run.Clean(ctx.Context, logger)
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update Plugins",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "plugin",
						Aliases: []string{"p"},
						Usage:   "Plugin to update",
					},
				},
				Action: func(ctx *cli.Context) error {
					logger = log.New(log.Level(ctx.String("level")), log.OutputText)
					plugins := ctx.StringSlice("plugin")
					return run.Update(ctx.Context, logger, plugins)
				},
			},
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Install Plugins",
				Action: func(ctx *cli.Context) error {
					logger = log.New(log.Level(ctx.String("level")), log.OutputText)
					return run.Install(ctx.Context, logger)
				},
			},
			{
				Name:    "source",
				Aliases: []string{"s"},
				Usage:   "Source Plugins",
				Action: func(ctx *cli.Context) error {
					logger = log.New(log.Level(ctx.String("level")), log.OutputText)
					return run.Source(ctx.Context, logger)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		if logger != nil {
			logger.Error("failed to run application", "error", err)
		} else {
			fmt.Println(err)
		}
	}
}
