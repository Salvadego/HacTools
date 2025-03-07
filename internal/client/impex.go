package client

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/url"

	"github.com/anaskhan96/soup"
	"github.com/matsal007/hactools/internal/logger"
)

func (c *HACClient) ImportScriptImpex(data map[string]any) (string, error) {
	logger.Info("Executing impex")
	logger.Debug("Impex data %+v", data)

	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, fmt.Sprintf("%v", value))
	}

	body, err := c.Post("console/impex/import", formData)
	if err != nil {
		return "", fmt.Errorf("failed to execute impex: %w", err)
	}

	logger.Debug("Response body: %s", string(body))

	doc := soup.HTMLParse(string(body))
	resultTag := doc.Find("div", "class", "impexResult")
	var result string
	if resultTag.Error != nil {
		result = ""
		return result, nil
	}

	result = resultTag.FullText()
	return result, nil
}

func (c *HACClient) ImportFileImpex(body *bytes.Buffer, writer *multipart.Writer) (string, error) {

	resp, err := c.PostMultipart("console/impex/import/upload", body, writer.FormDataContentType())
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	logger.Debug("Upload response body: %s", string(resp))

	doc := soup.HTMLParse(string(resp))
	resultTag := doc.Find("div", "class", "impexResult")

	var result string
	if resultTag.Error != nil {
		result = ""
		return result, nil
	}

	result = resultTag.FullText()
	return result, nil
}
