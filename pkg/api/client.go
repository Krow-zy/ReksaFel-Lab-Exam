package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// TelemetryClient dispatches local reports to the proctor dashboard.
type TelemetryClient struct {
	BaseURL     string
	SecretToken string
	HTTPClient  *http.Client
}

// NewTelemetryClient initializes a client with a standard 5-second timeout.
func NewTelemetryClient(baseURL, token string) *TelemetryClient {
	return &TelemetryClient{
		BaseURL:     baseURL,
		SecretToken: token,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Dispatch sends a JSON payload to the target endpoint with Bearer auth headers.
func (c *TelemetryClient) Dispatch(ctx context.Context, endpoint string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.SecretToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("dispatch failed with status: " + resp.Status)
	}

	return nil
}
