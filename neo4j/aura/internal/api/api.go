package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/neo4j/cli/common/clictx"
	"github.com/spf13/cobra"
)

const userAgent = "Neo4jCLI/%s"

type Grant struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func MakeRequest(cmd *cobra.Command, method string, path string, data map[string]any) (responseBody []byte, statusCode int, err error) {
	cmd.SilenceUsage = true

	client := http.Client{}
	var body io.Reader
	if data == nil {
		body = nil
	} else {
		jsonData, err := json.Marshal(data)

		if err != nil {
			return responseBody, 0, err
		}

		body = bytes.NewBuffer(jsonData)
	}

	config, ok := clictx.Config(cmd.Context())

	if !ok {
		return responseBody, 0, errors.New("error fetching cli configuration values")
	}

	baseUrl, err := config.GetString("aura.base-url")
	if err != nil {
		return responseBody, 0, err
	}

	u, _ := url.ParseRequestURI(baseUrl)
	u = u.JoinPath(path)
	urlString := u.String()

	req, err := http.NewRequest(method, urlString, body)

	if err != nil {
		return responseBody, 0, err
	}

	req.Header, err = getHeaders(cmd.Context())
	if err != nil {
		return responseBody, 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return responseBody, 0, err
	}

	defer res.Body.Close()

	if isSuccessful(res.StatusCode) {
		responseBody, err = io.ReadAll(res.Body)

		if err != nil {
			return responseBody, 0, err
		}

		return responseBody, res.StatusCode, nil
	}

	return responseBody, res.StatusCode, handleResponseError(res)
}

// Checks status code is 2xx
func isSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
