package main

import (
	"fmt"
	"os"

	"github.com/Salvadego/HacTools/internal/client"
	"github.com/Salvadego/HacTools/internal/flexsearch"
	"github.com/Salvadego/HacTools/internal/logger"
	"github.com/Salvadego/HacTools/internal/models"
	"github.com/Salvadego/HacTools/internal/options"
	"github.com/spf13/cobra"
)

var (
	maxCount    int
	noAnalyze   bool
	noBlacklist bool
	logLevel    string
)

var columnBlacklist = []string{
	"hjmpTS",
	"createdTS",
	"modifiedTS",
	"TypePkString",
	"OwnerPkString",
	"aCLTS",
	"propTS",
}

var conf options.Config

func init() {
	options.GetDefaults(rootCmd, &conf)
	rootCmd.PersistentFlags().IntVarP(&maxCount, "max-count", "m", 10, "Maximum number of results")
	rootCmd.PersistentFlags().BoolVarP(&noAnalyze, "no-analyze", "A", false, "Do not analyze PK")
	rootCmd.PersistentFlags().BoolVarP(&noBlacklist, "no-blacklist", "B", false, "Ignore column blacklist")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "error", "Log level (debug, info, error, none)")
}

var rootCmd = &cobra.Command{
	Use:          "xf [query or file path]",
	Short:        "Execute flexible search queries against Hybris HAC",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetLogLevel(logger.LogLevelFromString(logLevel))

		var query string
		arg := args[0]

		if _, err := os.Stat(arg); err == nil {
			data, err := os.ReadFile(arg)
			if err != nil {
				return fmt.Errorf("failed to read script file: %w", err)
			}
			query = string(data)
		} else {
			query = arg
		}

		if query == "" {
			return fmt.Errorf("query cannot be empty")
		}

		client := client.NewHACClient(conf.Address, conf.User, conf.Password)
		if err := client.Login(); err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}

		executor := flexsearch.NewFlexSearchExecutor(client)
		result, err := executor.Execute(query, models.FlexExecuteOptions{
			MaxCount:        maxCount,
			NoAnalyze:       noAnalyze,
			ColumnBlacklist: columnBlacklist,
			NoBlacklist:     noBlacklist,
		})

		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}

		return executor.DisplayResults(result)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
