package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/neo4j/cli/common/clicfg"
)

const userAgent = "Neo4jCLI/%s"

type Grant struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type RequestConfig struct {
	Method      string
	PostBody    map[string]any
	QueryParams map[string]string
}

func MakeRequest(cfg *clicfg.Config, path string, config *RequestConfig) (responseBody []byte, statusCode int, err error) {
	client := http.Client{}
	var method = config.Method
	if method == "" {
		panic(fmt.Sprintf("method not set in requests %s", path))
	}

	body, err := createBody(config.PostBody)
	if err != nil {
		return nil, 0, err
	}

	baseUrl := cfg.Aura.BaseUrl()

	u, _ := url.ParseRequestURI(baseUrl)
	u = u.JoinPath(path)

	addQueryParams(u, config.QueryParams)

	urlString := u.String()
	req, err := http.NewRequest(method, urlString, body)

	if err != nil {
		return responseBody, 0, err
	}

	credential, err := cfg.Aura.DefaultCredential()
	if err != nil {
		return responseBody, 0, err
	}

	req.Header, err = getHeaders(credential, cfg)
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

	return responseBody, res.StatusCode, handleResponseError(res, credential)
}

func createBody(data map[string]any) (io.Reader, error) {
	if data == nil {
		return nil, nil
	} else {
		jsonData, err := json.Marshal(data)

		if err != nil {
			return nil, err
		}

		return bytes.NewBuffer(jsonData), nil
	}
}

func addQueryParams(u *url.URL, params map[string]string) {
	if params != nil {
		q := u.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		u.RawQuery = q.Encode()
	}
}

// Checks status code is 2xx
func isSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
