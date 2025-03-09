package models

import "github.com/spf13/cobra"

type EditorConfig struct {
	FilePattern    string
	InitialContent string
	ExecutorFunc   func(string) error
	CustomFlags    []func(*cobra.Command)
}
