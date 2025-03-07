package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/SalvadegoDev/HacTools/internal/logger"
	"github.com/SalvadegoDev/HacTools/internal/models"
)

func (c *HACClient) ExecuteFlexSearch(data map[string]any, blacklist []string) (*models.FlexSearchResponse, error) {
	logger.Info("Executing flex search")
	logger.Debug("Query data: %+v", data)

	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, fmt.Sprintf("%v", value))
	}

	body, err := c.Post("console/flexsearch/execute", formData)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Debug("Response body: %s", string(body))

	var result models.FlexSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	if len(result.ResultList) > 0 {
		validColumns := make([]int, 0)
		for colIdx, header := range result.Headers {
			if isBlacklisted(header, blacklist) {
				continue
			}

			hasValue := false
			for rowIdx := range result.ResultList {
				if result.ResultList[rowIdx][colIdx] != "" && result.ResultList[rowIdx][colIdx] != "null" {
					hasValue = true
					break
				}
			}
			if hasValue {
				validColumns = append(validColumns, colIdx)
			}
		}

		if len(validColumns) < len(result.Headers) {
			newHeaders := make([]string, len(validColumns))
			for newIdx, oldIdx := range validColumns {
				newHeaders[newIdx] = result.Headers[oldIdx]
			}
			result.Headers = newHeaders

			newResultList := make([][]string, len(result.ResultList))
			for rowIdx, row := range result.ResultList {
				newRow := make([]string, len(validColumns))
				for newIdx, oldIdx := range validColumns {
					newRow[newIdx] = row[oldIdx]
				}
				newResultList[rowIdx] = newRow
			}
			result.ResultList = newResultList
		}
	}

	return &result, nil
}

func isBlacklisted(header string, blacklist []string) bool {
	normalizedHeader := strings.ToLower(strings.TrimPrefix(strings.TrimPrefix(header, "p_"), "P_"))
	for _, blacklisted := range blacklist {
		if strings.ToLower(blacklisted) == normalizedHeader {
			return true
		}
	}
	return false
}
