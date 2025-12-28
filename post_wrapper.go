package ptti

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// TODO: can be better
func PostJSON[T any](url string, input any, result *T) error {
	var full struct {
		Response T `json:"response"`
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&full); err != nil {
		return err
	}

	*result = full.Response
	return nil
}

// TODO: can be better
func FetchJSON[T any](url string, result *T) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var full struct {
		Response T `json:"response"`
		Data     T `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&full); err != nil {
		return err
	}

	*result = full.Data
	return nil
}
