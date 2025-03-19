package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/Salvadego/HacTools/internal/logger"
	"github.com/Salvadego/HacTools/models"
)

func (c *HACClient) AnalyzePK(pk string) (*models.PKAnalyzeResponse, error) {
	logger.Info("Executing pk analyze")
	logger.Debug("PK: %+v", pk)

	formData := url.Values{}
	formData.Set("pkString", pk)
	formData.Set("_csrf", c.Csrf)

	body, err := c.Post("platform/pkanalyzer/analyze", formData)
	if err != nil {

		if strings.Contains(err.Error(), "HTTP 403") {
			logger.Info("Session appears to be expired, attempting to re-login")

			loginErr := c.Login()
			if loginErr != nil {
				return nil, fmt.Errorf("failed to re-login: %w", loginErr)
			}

			formData.Set("_csrf", c.Csrf)

			body, err = c.Post("platform/pkanalyzer/analyze", formData)
			if err != nil {
				return nil, fmt.Errorf("failed to analyze PK even after re-login: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to analyze PK: %w", err)
		}
	}

	logger.Debug("Response body from PK analyze: %s", string(body))

	var result models.PKAnalyzeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	return &result, nil
}
