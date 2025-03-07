package groovy

import (
	"fmt"

	"github.com/matsal007/hactools/internal/client"
	"github.com/matsal007/hactools/internal/models"
)

type GroovyExecutor struct {
	Client *client.HACClient
}

func NewGroovyExecutor(client *client.HACClient) *GroovyExecutor {
	return &GroovyExecutor{
		Client: client,
	}
}

func (e *GroovyExecutor) Execute(script string, opts models.GroovyExecuteOptions) (*models.GroovyResponse, error) {
	data := map[string]any{
		"script":     script,
		"_csrf":      e.Client.Csrf,
		"scriptType": opts.ScriptType,
		"commit":     opts.Commit,
	}

	resp, err := e.Client.ExecuteGroovy(data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *GroovyExecutor) DisplayResults(result *models.GroovyResponse) error {
	if result == nil {
		return fmt.Errorf("no results to display")
	}

	fmt.Println("=== OUTPUT ===")
	if result.ScriptResult != "" {
		fmt.Println(result.ScriptResult)
	} else {
		fmt.Println("<No return value>")
	}

	if result.ExecutionResult != "" {
		fmt.Println("\n=== RESULT ===")
		fmt.Println(result.ExecutionResult)
	}

	if result.StacktraceText != "" {
		fmt.Println("\n=== STACKTRACE ===")
		fmt.Println(result.StacktraceText)
		return fmt.Errorf("script execution failed with error")
	}

	return nil
}
