package implementations

import (
	"net/http"
	"strings"
	"time"

	"bos/constants"
)

const (
	beginSimulationPath   = "/v1/simulation/begin"
	performSimulationPath = "/v1/simulation/perform"
)

type Simulator struct {
	baseURL string
	client  *http.Client
}

func NewSimulator() *Simulator {
	return NewSimulatorWithBaseURL(constants.SimulatorBaseURL)
}

func NewSimulatorWithBaseURL(baseURL string) *Simulator {
	return NewSimulatorWithHTTPClient(baseURL, &http.Client{Timeout: 60 * time.Second})
}

func NewSimulatorWithHTTPClient(baseURL string, client *http.Client) *Simulator {
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	return &Simulator{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}
