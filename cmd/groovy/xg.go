package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Salvadego/HacTools/internal/client"
	"github.com/Salvadego/HacTools/internal/editor"
	"github.com/Salvadego/HacTools/internal/groovy"
	"github.com/Salvadego/HacTools/internal/logger"
	"github.com/Salvadego/HacTools/internal/options"
	"github.com/Salvadego/HacTools/models"
	"github.com/spf13/cobra"
)

var (
	commit            bool
	scriptType        string
	logLevel          string
	scriptFilePattern = map[string]string{
		"groovy":     "groovy-script-*.groovy",
		"javascript": "js-script-*.js",
		"branshell":  "branshell-script-*.branshell",
	}
)

var conf options.Config

func init() {
	options.GetDefaults(rootCmd, &conf)
	rootCmd.PersistentFlags().BoolVarP(&commit, "commit", "c", false, "Execute with commit")
	rootCmd.PersistentFlags().StringVarP(&scriptType, "type", "t", "groovy", "Script type (groovy, javascript, beanshell)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "error", "Log level (debug, info, error, none)")

	editorCommand := editor.CreateEditorCommand(models.EditorConfig{
		FilePattern:    scriptFilePattern[strings.ToLower(scriptType)],
		InitialContent: "// Enter your Groovy script here\n\n",
		ExecutorFunc: executorFunc,
		CustomFlags: []func(*cobra.Command){},
	})

	rootCmd.AddCommand(editorCommand)
}

var rootCmd = &cobra.Command{
	Use:          "xg [script or file path]",
	Short:        "Execute Groovy scripts against Hybris HAC",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SetLogLevel(logger.LogLevelFromString(logLevel))
		var script string
		arg := args[0]

		if _, err := os.Stat(arg); err == nil {
			data, err := os.ReadFile(arg)
			if err != nil {
				return fmt.Errorf("failed to read script file: %w", err)
			}
			script = string(data)
		} else {
			script = arg
		}

        return executorFunc(script)
	},
}

func executorFunc(script string) error {
	scriptType = strings.ToLower(scriptType)
	if scriptType != "groovy" && scriptType != "javascript" && scriptType != "beanshell" {
		return fmt.Errorf("invalid script type: %s (must be groovy, javascript, or beanshell)", scriptType)
	}

	client := client.NewHACClient(conf.Address, conf.User, conf.Password)
	if err := client.Login(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	executor := groovy.NewGroovyExecutor(client)
	result, err := executor.Execute(script, models.GroovyExecuteOptions{
		ScriptType: scriptType,
		Commit:     commit,
	})

	if err != nil {
		return fmt.Errorf("failed to execute script: %w", err)
	}

	return executor.DisplayResults(result)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
