package main

import (
	"fmt"
	"os"

	"github.com/Salvadego/HacTools/internal/client"
	"github.com/Salvadego/HacTools/internal/impex"
	"github.com/Salvadego/HacTools/internal/logger"
	"github.com/Salvadego/HacTools/models"
	"github.com/Salvadego/HacTools/internal/options"
	"github.com/spf13/cobra"
)

var (
	legacyMode          bool
	enableCodeExecution bool
	distributedMode     bool
	sldEnabled          bool
	logLevel            string
)

var conf options.Config

func init() {
	options.GetDefaults(rootCmd, &conf)
	rootCmd.PersistentFlags().BoolVarP(&legacyMode, "legacy", "L", false, "Enable legacyMode")
	rootCmd.PersistentFlags().BoolVarP(&enableCodeExecution, "exec", "c", false, "Enable code execution")
	rootCmd.PersistentFlags().BoolVarP(&distributedMode, "distributed", "d", false, "Enable legacyMode")
	rootCmd.PersistentFlags().BoolVarP(&sldEnabled, "sld", "", false, "Enable Sld")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "error", "Log level (debug, info, error, none)")
}

var rootCmd = &cobra.Command{
	Use:          "ii [script or file path]",
	Short:        "Import Impex against Hybris HAC",
	Long:         `A impex importer for Hybris HAC`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetLogLevel(logger.LogLevelFromString(logLevel))
		client := client.NewHACClient(conf.Address, conf.User, conf.Password)
		if err := client.Login(); err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}

		arg := args[0]

		importer := impex.NewImpexImporter(client)
		options := models.ImpexExecuteOptions{
			LegacyMode:          legacyMode,
			EnableCodeExecution: enableCodeExecution,
			DistributedMode:     distributedMode,
			SldEnabled:          sldEnabled,
		}

		if _, err := os.Stat(arg); err == nil {
			result, err := importer.ImportFile(arg, options)

			if err != nil {
				return fmt.Errorf("failed to execute script: %w", err)
			}

			return importer.DisplayResults(result)
		}

		var script string
		script = arg
		result, err := importer.ImportScript(script, options)
		if err != nil {
			return fmt.Errorf("failed to execute script: %w", err)
		}

		return importer.DisplayResults(result)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
