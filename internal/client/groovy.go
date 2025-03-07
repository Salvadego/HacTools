package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SalvadegoDev/HacTools/internal/logger"
	"github.com/SalvadegoDev/HacTools/internal/models"
)

func (c *HACClient) ExecuteGroovy(data map[string]any) (*models.GroovyResponse, error) {
	logger.Info("Executing script")
	logger.Debug("Script data: %+v", data)

	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, fmt.Sprintf("%v", value))
	}

	body, err := c.Post("console/scripting/execute", formData)
	if err != nil {
		return nil, fmt.Errorf("failed to execute script: %w", err)
	}

	logger.Debug("Response body: %s", string(body))

	var result models.GroovyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	return &result, nil
}
