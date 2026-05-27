package implementations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"bos/constants"
)

func (s *Simulator) postJSON(path string, payload any, target any) error {
	responseBody, err := s.postJSONRaw(path, payload)
	if err != nil {
		return err
	}
	if len(responseBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(responseBody, target); err != nil {
		return fmt.Errorf("decode simulator API response %s: %w", path, err)
	}
	return nil
}

func (s *Simulator) postJSONRaw(path string, payload any) ([]byte, error) {
	if s.client == nil {
		s.client = &http.Client{Timeout: 60 * time.Second}
	}
	if s.baseURL == "" {
		s.baseURL = constants.SimulatorBaseURL
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("simulator API %s: %w", path, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read simulator API response %s: %w", path, err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("simulator API %s returned %s: %s", path, resp.Status, strings.TrimSpace(string(responseBody)))
	}
	return responseBody, nil
}
