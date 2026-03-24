package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const clientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

type tokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func exchangeToken(opts options, assertion string) (string, error) {
	form := url.Values{
		"grant_type":            {"client_credentials"},
		"client_assertion_type": {clientAssertionType},
		"client_assertion":      {assertion},
		"scope":                 {opts.Scope},
	}

	req, err := http.NewRequest(http.MethodPost, opts.Audience, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var payload tokenExchangeResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("unexpected response status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if payload.Error != "" {
			return "", fmt.Errorf("%s: %s", payload.Error, payload.Description)
		}
		return "", fmt.Errorf("token endpoint returned %s", resp.Status)
	}

	if payload.AccessToken == "" {
		return "", fmt.Errorf("token endpoint response did not include access_token")
	}

	return payload.AccessToken, nil
}
