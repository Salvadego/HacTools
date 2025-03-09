package models

type GroovyExecuteOptions struct {
	ScriptType string
	Commit     bool
}

type GroovyResponse struct {
	ExecutionResult string `json:"executionResult"`
	ScriptResult    string `json:"outputText"`
	StacktraceText  string `json:"stacktraceText"`
	ExceptionText   string `json:"exceptionText"`
	Success         bool   `json:"success"`
}
