package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clicfg/credentials"
)

func getToken(credential *credentials.AuraCredential, cfg *clicfg.Config) (string, error) {
	if credential.HasValidAccessToken() {
		return credential.AccessToken, nil
	}

	data := url.Values{}

	data.Set("grant_type", "client_credentials")

	url := cfg.Aura.AuthUrl()

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	version := cfg.Version

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
	defer res.Body.Close()

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

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	var grant Grant

	err = json.Unmarshal(resBody, &grant)

	if err != nil {
		return "", err
	}

	_, err = cfg.Credentials.Aura.UpdateAccessToken(credential, grant.AccessToken, grant.ExpiresIn)
	return grant.AccessToken, err
}
