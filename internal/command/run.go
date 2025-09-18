package command

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kukymbr/dbmodgen/internal/generator"
	"github.com/kukymbr/dbmodgen/internal/util"
	"github.com/kukymbr/dbmodgen/internal/version"
	"github.com/spf13/cobra"
)

func Run() {
	if err := run(); err != nil {
		util.PrintErrorf("%s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	var (
		configPath string
		silent     bool
	)

	var cmd = &cobra.Command{
		Use:   "dbmodgen",
		Short: "DB models generator",
		Long:  `Generates models from the existing database`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			opt, err := generator.ReadOptions(configPath, os.Getenv)
			if err != nil {
				return err
			}

			gen, err := generator.New(opt)
			if err != nil {
				return err
			}

			return gen.Generate(ctx)
		},
		Version: version.GetVersion(),
	}

	initFlags(cmd, &configPath, &silent)

	cmd.PersistentPreRun = func(_ *cobra.Command, _ []string) {
		util.SetSilentMode(silent)
	}

	return cmd.Execute()
}

func initFlags(cmd *cobra.Command, configPath *string, silent *bool) {
	cmd.PersistentFlags().BoolVarP(silent, "silent", "s", false, "Silent mode")

	cmd.Flags().StringVarP(
		configPath,
		"config",
		"c",
		"",
		"Target package name of the generated code",
	)

	_ = cmd.MarkFlagRequired("conf")
	_ = cmd.MarkFlagFilename("conf")
}
