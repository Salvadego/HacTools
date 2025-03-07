package impex

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/SalvadegoDev/HacTools/internal/client"
	"github.com/SalvadegoDev/HacTools/internal/models"
)

type ImpexImporter struct {
	Client *client.HACClient
}

func NewImpexImporter(client *client.HACClient) *ImpexImporter {
	return &ImpexImporter{
		Client: client,
	}
}

func boolToOnOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

func (e *ImpexImporter) ImportScript(script string, opts models.ImpexExecuteOptions) (string, error) {

	data := map[string]any{
		"_csrf":                e.Client.Csrf,
		"scriptContent":        script,
		"validationEnum":       "IMPORT_STRICT",
		"maxThreads":           "16",
		"encoding":             "UTF-8",
		"_enableCodeExecution": boolToOnOff(opts.EnableCodeExecution),
		"_distributedMode":     boolToOnOff(opts.DistributedMode),
		"_legacyMode":          boolToOnOff(opts.LegacyMode),
		"_sldEnabled":          boolToOnOff(opts.SldEnabled),
	}

	resp, err := e.Client.ImportScriptImpex(data)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (e *ImpexImporter) ImportFile(filepath string, opts models.ImpexExecuteOptions) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	formFields := map[string]string{
		"_csrf":               e.Client.Csrf,
		"maxThreads":          "16",
		"validationEnum":      "IMPORT_STRICT",
		"encoding":            "UTF-8",
		"enableCodeExecution": boolToOnOff(opts.EnableCodeExecution),
		"distributedMode":     boolToOnOff(opts.DistributedMode),
		"legacyMode":          boolToOnOff(opts.LegacyMode),
		"sldEnabled":          boolToOnOff(opts.SldEnabled),
	}

	for key, val := range formFields {
		err = writer.WriteField(key, val)
		if err != nil {
			return "", fmt.Errorf("failed to write form field: %w", err)
		}
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	result, err := e.Client.ImportFileImpex(body, writer)

	if err != nil {
		return "", err
	}
	return result, nil
}

func (e *ImpexImporter) DisplayResults(result string) error {
	if result != "" {
		return fmt.Errorf("%s", result)
	}

	fmt.Println("=== RESULT ===")
	fmt.Println("Import finished successfully")

	return nil
}
