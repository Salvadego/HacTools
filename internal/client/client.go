package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/Salvadego/HacTools/internal/logger"
)

type HACClient struct {
	Client   *http.Client
	BaseURL  string
	Username string
	Password string
	Csrf     string
}

func NewHACClient(baseURL, username, password string) *HACClient {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return &HACClient{
		Client:   client,
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}
}

func (c *HACClient) clearSession() {
	logger.Debug("Clearing session and CSRF token")
	c.Client.Jar, _ = cookiejar.New(nil)
	c.Csrf = ""
}

func (c *HACClient) extractCSRFToken(body string) (string, error) {
	re := regexp.MustCompile(`name="_csrf"\s+value="(.+?)"\s*/>`)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("CSRF token not found in response")
	}
	return matches[1], nil
}

func (c *HACClient) validateLoginResponse(body string) bool {
	return strings.Contains(body, "You're")
}

func (c *HACClient) getInitialCSRFToken() (string, error) {
	logger.Info("Getting initial CSRF token")
	resp, err := c.Client.Get(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("failed to get login page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login page: %w", err)
	}

	if strings.Contains(string(body), "503: This service is down for maintenance") ||
		strings.Contains(string(body), "SAP Commerce Cloud - Maintenance") {
		return "", fmt.Errorf("service is down for maintenance")
	}

	csrf, err := c.extractCSRFToken(string(body))
	if err != nil {
		return "", fmt.Errorf("failed to extract CSRF token: %w", err)
	}

	return csrf, nil
}

func (c *HACClient) Login() error {
	logger.Info("Starting login process")
	c.clearSession()

	initialCSRF, err := c.getInitialCSRFToken()
	if err != nil {
		return fmt.Errorf("failed to get initial CSRF token: %w", err)
	}
	logger.Debug("Initial CSRF token: %s", initialCSRF)

	loginData := url.Values{
		"j_username":                   {c.Username},
		"j_password":                   {c.Password},
		"_csrf":                        {initialCSRF},
		"_spring_security_remember_me": {"true"},
	}

	loginURL := c.BaseURL + "j_spring_security_check"
	logger.Debug("Sending credentials to %s", loginURL)

	resp, err := c.Client.PostForm(loginURL, loginData)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	if !c.validateLoginResponse(string(body)) {
		resp, err = c.Client.Get(c.BaseURL)
		if err != nil {
			return fmt.Errorf("failed to verify login: %w", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read verification page: %w", err)
		}

		if !c.validateLoginResponse(string(body)) {
			return fmt.Errorf("login failed: invalid credentials")
		}
	}

	csrf, err := c.extractCSRFToken(string(body))
	if err != nil {
		return fmt.Errorf("failed to get post-login CSRF token: %w", err)
	}
	c.Csrf = csrf
	logger.Info("Successfully logged in with CSRF token: %s", c.Csrf)

	return nil
}

func (c *HACClient) Post(endpoint string, data url.Values) ([]byte, error) {
	resp, err := c.Client.PostForm(c.BaseURL+endpoint, data)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received HTTP %d error: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == 405 {
		c.clearSession()
		c.getInitialCSRFToken()
		return c.Post(endpoint, data)
		// return nil, fmt.Errorf("received HTTP %d error: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *HACClient) PostMultipart(endpoint string, body *bytes.Buffer, contentType string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}
