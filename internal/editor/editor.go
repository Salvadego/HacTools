package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Salvadego/HacTools/models"
	"github.com/spf13/cobra"
)

func OpenEditor(initialContent string, filePattern string) (string, error) {
	editorCmd := os.Getenv("EDITOR")
	if editorCmd == "" {
		editorCmd = "vi"
	}

	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, filePattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if initialContent != "" {
		if _, err := tempFile.WriteString(initialContent); err != nil {
			return "", fmt.Errorf("failed to write to temporary file: %w", err)
		}
	}
	
	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file: %w", err)
	}

	cmdParts := strings.Fields(editorCmd)
	cmd := exec.Command(cmdParts[0], append(cmdParts[1:], tempFile.Name())...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read edited file: %w", err)
	}

	return string(content), nil
}

func CreateEditorCommand(opts models.EditorConfig) *cobra.Command {
	var savePath string
	
	cmd := &cobra.Command{
		Use:   "editor [template_file]",
		Short: "Open content in editor",
		Long: `Opens your preferred text editor to create or edit content.
The editor is determined by the $EDITOR environment variable. If not set, vi will be used.
You can optionally provide a template file as an argument.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			initialContent := opts.InitialContent
			if len(args) > 0 {
				templatePath := args[0]
				content, err := os.ReadFile(templatePath)
				if err != nil {
					return fmt.Errorf("failed to read template file: %w", err)
				}
				initialContent = string(content)
			}

			content, err := OpenEditor(initialContent, opts.FilePattern)
			if err != nil {
				return fmt.Errorf("editor error: %w", err)
			}

			if content == "" {
				return fmt.Errorf("content cannot be empty")
			}

			if savePath != "" {
				saveDir := filepath.Dir(savePath)
				if _, err := os.Stat(saveDir); os.IsNotExist(err) {
					if err := os.MkdirAll(saveDir, 0755); err != nil {
						return fmt.Errorf("failed to create directory for saving: %w", err)
					}
				}
				
				if err := os.WriteFile(savePath, []byte(content), 0644); err != nil {
					return fmt.Errorf("failed to save content: %w", err)
				}
				fmt.Printf("Content saved to %s\n", savePath)
			}

			return opts.ExecutorFunc(content)
		},
	}

	cmd.Flags().StringVar(&savePath, "save", "", "Save content to file after editing")

	for _, addFlag := range opts.CustomFlags {
		addFlag(cmd)
	}

	return cmd
}

