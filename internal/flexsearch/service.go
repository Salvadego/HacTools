package flexsearch

import (
	"fmt"
	"html"
	"os"
	"os/exec"
	"strings"

	"github.com/Salvadego/HacTools/internal/client"
	"github.com/Salvadego/HacTools/models"
	"github.com/olekukonko/tablewriter"
)

type FlexSearchExecutor struct {
	Client *client.HACClient
}

func NewFlexSearchExecutor(client *client.HACClient) *FlexSearchExecutor {
	executor := &FlexSearchExecutor{
		Client: client,
	}
	return executor
}

func (e *FlexSearchExecutor) Execute(query string, opts models.FlexExecuteOptions) (*models.FlexSearchResponse, error) {
	data := map[string]any{
		"flexibleSearchQuery": query,
		"_csrf":               e.Client.Csrf,
		"maxCount":            opts.MaxCount,
		"user":                e.Client.Username,
		"locale":              "en",
		"commit":              false,
	}

	var blacklist []string
	if !opts.NoBlacklist {
		blacklist = opts.ColumnBlacklist
	}
	resp, err := e.Client.ExecuteFlexSearch(data, blacklist)
	if err != nil {
		return nil, err
	}

	if resp.Exception != nil {
		return nil, fmt.Errorf("flex search error: %s", resp.Exception.Message)
	}

	return resp, nil
}

func (e *FlexSearchExecutor) formatTable(result *models.FlexSearchResponse) string {
	var buf strings.Builder
	table := tablewriter.NewWriter(&buf)

	table.SetBorder(false)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetAutoWrapText(true)
	table.SetRowLine(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	headers := make([]string, len(result.Headers))
	for i, h := range result.Headers {
		headers[i] = strings.TrimPrefix(strings.TrimPrefix(h, "p_"), "P_")
	}
	table.SetHeader(headers)

	table.SetColWidth(20)

	for _, row := range result.ResultList {
		tableRow := make([]string, len(row))
		for i, cell := range row {
			if cell != "" {
				tableRow[i] = html.UnescapeString(cell)
			} else {
				tableRow[i] = ""
			}
		}
		table.Append(tableRow)
	}

	footerRow := make([]string, len(headers))
	footerRow[0] = fmt.Sprintf("Total Rows: %d", len(result.ResultList))
	table.SetFooter(footerRow)

	table.Render()
	return buf.String()
}

func (e *FlexSearchExecutor) DisplayResults(result *models.FlexSearchResponse) error {
	if result == nil {
		return fmt.Errorf("no results to display")
	}

	if len(result.ResultList) == 0 {
		fmt.Println("No results found")
		return nil
	}

	tableOutput := e.formatTable(result)

	if isPipe() {
		fmt.Print(tableOutput)
		return nil
	}

	return e.displayWithPager(tableOutput)
}

func (e *FlexSearchExecutor) displayWithPager(content string) error {
	tmpfile, err := os.CreateTemp("", "flexsearch-*.txt")
	fmt.Println(tmpfile.Name())
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := tmpfile.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		pagerCmd = "less -RS"
	}
	args := strings.Fields(pagerCmd)
	args = append(args, tmpfile.Name())

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func isPipe() bool {
	fi, _ := os.Stdout.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}
