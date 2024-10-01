package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/neo4j/cli/common/clicfg"
)

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Field   string `json:"field"`
}

type ServerError struct {
	Error string `json:"error"`
}

func handleResponseError(res *http.Response) error {
	var err error
	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	switch statusCode := res.StatusCode; statusCode {
	// redirection messages
	case http.StatusPermanentRedirect:
		return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
	// client error responses
	case http.StatusBadRequest:
		var errorResponse ErrorResponse

		err = json.Unmarshal(resBody, &errorResponse)
		if err != nil {
			return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
		}

		messages := []string{}
		for _, e := range errorResponse.Errors {
			messages = append(messages, e.Message)
		}

		return fmt.Errorf("%s", messages)
	case http.StatusUnauthorized:
		// TODO: clear the token?
		var errorResponse ErrorResponse

		err = json.Unmarshal(resBody, &errorResponse)
		if err != nil {
			return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
		}

		messages := []string{}
		for _, e := range errorResponse.Errors {
			messages = append(messages, e.Message)
		}

		return fmt.Errorf("%s", messages)
	case http.StatusForbidden:
		var serverError ServerError

		err := json.Unmarshal(resBody, &serverError)

		if err != nil {
			return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
		}

		// TODO: clear the token?
		var errorResponse ErrorResponse

		err = json.Unmarshal(resBody, &errorResponse)
		if err != nil {
			return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
		}

		messages := []string{}
		for _, e := range errorResponse.Errors {
			messages = append(messages, e.Message)
		}

		return fmt.Errorf("%s", messages)
	case http.StatusNotFound:
		var errorResponse ErrorResponse

		err = json.Unmarshal(resBody, &errorResponse)
		if err != nil {
			return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
		}

		messages := []string{}
		for _, e := range errorResponse.Errors {
			messages = append(messages, e.Message)
		}

		return fmt.Errorf("%s", messages)
	case http.StatusMethodNotAllowed:
		var errorResponse ErrorResponse

		err = json.Unmarshal(resBody, &errorResponse)
		if err != nil {
			return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
		}

		messages := []string{}
		for _, e := range errorResponse.Errors {
			messages = append(messages, e.Message)
		}

		return fmt.Errorf("%s", messages)
	case http.StatusConflict:
		var errorResponse ErrorResponse

		err = json.Unmarshal(resBody, &errorResponse)
		if err != nil {
			return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
		}

		messages := []string{}
		for _, e := range errorResponse.Errors {
			messages = append(messages, e.Message)
		}

		return fmt.Errorf("%s", messages)
	case http.StatusUnsupportedMediaType:
		return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
	case http.StatusTooManyRequests:
		retryAfter := res.Header.Get("Retry-After")
		return fmt.Errorf("server rate limit exceeded, suggested cool-off period is %s seconds before rerunning the command", retryAfter)
	// server error responses
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		var errorResponse ErrorResponse

		err = json.Unmarshal(resBody, &errorResponse)
		if err != nil {
			return fmt.Errorf("unexpected error running CLI with args %s, please report an issue in https://github.com/neo4j/cli", os.Args[1:])
		}

		messages := []string{}
		for _, e := range errorResponse.Errors {
			messages = append(messages, e.Message)
		}

		return fmt.Errorf("%s", messages)
	default:
		return fmt.Errorf("unexpected status code %d and body %s running CLI with args %s, please report an issue in https://github.com/neo4j/cli", statusCode, resBody, os.Args[1:])
	}
}

func getHeaders(cfg *clicfg.Config) (http.Header, error) {
	token, err := getToken(cfg)

	if err != nil {
		return nil, err
	}

	version := cfg.Version

	return http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
		"User-Agent":    {fmt.Sprintf(userAgent, version)},
	}, nil
}

// Response types

const (
	InstanceStatusCreating      string = "creating"
	InstanceStatusDestroying    string = "destroying"
	InstanceStatusRunning       string = "running"
	InstanceStatusPausing       string = "pausing"
	InstanceStatusPaused        string = "paused"
	InstanceStatusSuspending    string = "suspending"
	InstanceStatusSuspended     string = "suspended"
	InstanceStatusResuming      string = "resuming"
	InstanceStatusLoading       string = "loading"
	InstanceStatusLoadingFailed string = "loading failed"
	InstanceStatusRestoring     string = "restoring"
	InstanceStatusUpdating      string = "updating"
	InstanceStatusOverwriting   string = "overwriting"
)

const (
	SnapshotStatusPending    string = "Pending"
	SnapshotStatusCompleted  string = "Completed"
	SnapshotStatusInProgress string = "InProgress"
	SnapshotStatusFailed     string = "Failed"
)

// Response Body of Create and Get Instance for successful requests
type CreateInstanceResponse struct {
	Data struct {
		Id            string
		ConnectionUrl string `json:"connection_url"`
		Username      string
		Password      string
		TenantId      string `json:"tenant_id"`
		CloudProvider string `json:"cloud_provider"`
		Region        string
		Type          string
		Name          string
	}
}

const (
	CMKStatusReady   = "ready"
	CMKStatusPending = "pending"
)

// Response Body of Create and Get Instance for successful requests
type CreateCMKResponse struct {
	Data struct {
		Id     string
		Status string
	}
}

// Response Body of Create and Get Instance for successful requests
type CreateSnapshotResponse struct {
	Data struct {
		SnapshotId string `json:"snapshot_id"`
	}
}

func ParseBody(body []byte) ([]map[string]any, error) {
	var values []map[string]any
	var jsonWithArray struct{ Data []map[string]any }

	err := json.Unmarshal(body, &jsonWithArray)

	// Try unmarshalling array first, if not it creates an array from the single item
	if err == nil {
		values = jsonWithArray.Data
	} else {
		var jsonWithSingleItem struct{ Data map[string]any }
		err := json.Unmarshal(body, &jsonWithSingleItem)
		if err != nil {
			return nil, err
		}
		values = []map[string]any{jsonWithSingleItem.Data}
	}

	return values, nil
}
