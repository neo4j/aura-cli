package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/neo4j/cli/common/clictx"
	"github.com/spf13/cobra"
)

const userAgent = "Neo4jCLI/%s"

type Grant struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

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

func getToken(ctx context.Context) (string, error) {
	config, ok := clictx.Config(ctx)
	if !ok {
		return "", errors.New("error fetching cli configuration values")
	}

	credential, err := config.Aura.GetDefaultCredential()
	if err != nil {
		return "", err
	}

	if credential.IsAccessTokenValid() {
		return credential.AccessToken, nil
	}

	data := url.Values{}

	data.Set("grant_type", "client_credentials")

	u, err := config.Get("aura.auth-url")

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", u.(string), strings.NewReader(data.Encode()))

	if err != nil {
		return "", err
	}

	version, ok := clictx.Version(ctx)
	if !ok {
		return "", errors.New("error fetching version from context")
	}

	req.Header = http.Header{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"User-Agent":   {fmt.Sprintf(userAgent, version)},
	}
	req.SetBasicAuth(credential.ClientId, credential.ClientSecret)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	switch statusCode := res.StatusCode; statusCode {
	case http.StatusBadRequest:
		return "", errors.New("request is invalid")
	case http.StatusUnauthorized:
		return "", errors.New("the provided credentials are invalid, expired, or revoked")
	case http.StatusForbidden:
		return "", errors.New("the request body is invalid")
	case http.StatusNotFound:
		return "", errors.New("the request body is missing")
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	var grant Grant

	err = json.Unmarshal(resBody, &grant)

	if err != nil {
		return "", err
	}

	credential.UpdateAccessToken(grant.AccessToken, grant.ExpiresIn)
	config.Write()
	return grant.AccessToken, nil
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

func MakeRequest(cmd *cobra.Command, method string, path string, data map[string]any) error {
	cmd.SilenceUsage = true

	client := http.Client{}
	var body io.Reader
	if data == nil {
		body = nil
	} else {
		jsonData, err := json.Marshal(data)

		if err != nil {
			return err
		}

		body = bytes.NewBuffer(jsonData)
	}

	config, ok := clictx.Config(cmd.Context())

	if !ok {
		return errors.New("error fetching cli configuration values")
	}

	baseUrl, err := config.Get("aura.base-url")
	if err != nil {
		return err
	}

	u, _ := url.ParseRequestURI(baseUrl.(string))
	u = u.JoinPath(path)
	urlString := u.String()

	req, err := http.NewRequest(method, urlString, body)

	if err != nil {
		return err
	}

	req.Header, err = getHeaders(cmd.Context())

	if err != nil {
		return err
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	output, err := config.Get("aura.output")
	if err != nil {
		return err
	}

	switch statusCode := res.StatusCode; statusCode {
	// successful responses
	case http.StatusOK:
		fallthrough
	case http.StatusAccepted:
		if len(resBody) > 0 {
			switch output := output.(string); output {
			case "json":
				var pretty bytes.Buffer
				err := json.Indent(&pretty, resBody, "", "\t")
				if err != nil {
					return err
				}
				cmd.Println(pretty.String())
			default:
				cmd.Println(string(resBody))
			}
		}
	case http.StatusNoContent:
		cmd.Println("Operation Successful")
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
	case http.StatusInternalServerError:
		fallthrough
	case http.StatusBadGateway:
		fallthrough
	case http.StatusServiceUnavailable:
		fallthrough
	case http.StatusGatewayTimeout:
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

	return nil
}
