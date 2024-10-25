package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HealthCheckRoute struct {
	EnableRouteHealthCheckError bool
	Routes                      []Route
}

// HealthCheckResponse represents the health check response structure
type HealthCheckResponse struct {
	Status string                     `json:"status"`
	Routes []HealthCheckRouteResponse `json:"routes"`
}
type HealthCheckRouteResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func HealthCheck(healthURL string) error {
	healthCheckURL, err := url.Parse(healthURL)
	if err != nil {
		return fmt.Errorf("error parsing HealthCheck URL: %v ", err)
	}
	// Create a new request for the route
	healthReq, err := http.NewRequest("GET", healthCheckURL.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating HealthCheck request: %v ", err)
	}
	// Perform the request to the route's healthcheck
	client := &http.Client{}
	healthResp, err := client.Do(healthReq)
	if err != nil {
		return fmt.Errorf("error performing HealthCheck request: %v ", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(healthResp.Body)

	if healthResp.StatusCode >= 400 {
		return fmt.Errorf("health check failed with status code %v", healthResp.StatusCode)
	}
	return nil
}
