package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/neo4j/cli/common/clictx"
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

func getHeaders(ctx context.Context) (http.Header, error) {
	token, err := getToken(ctx)

	if err != nil {
		return nil, err
	}

	version, ok := clictx.Version(ctx)
	if !ok {
		return nil, errors.New("error fetching version from context")
	}

	return http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
		"User-Agent":    {fmt.Sprintf(userAgent, version)},
	}, nil
}
