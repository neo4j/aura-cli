package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/neo4j/cli/common/clicfg"
)

const userAgent = "Neo4jCLI/%s"

type Grant struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func MakeRequest(cfg *clicfg.Config, method string, path string, data map[string]any) (responseBody []byte, statusCode int, err error) {
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

	baseUrl := cfg.Aura.BaseUrl()

	u, _ := url.ParseRequestURI(baseUrl)
	u = u.JoinPath(path)
	urlString := u.String()

	// Quick fix
	urlString = strings.Replace(urlString, "%3F", "?", 1)
	req, err := http.NewRequest(method, urlString, body)

	if err != nil {
		return responseBody, 0, err
	}

	req.Header, err = getHeaders(cfg)
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
