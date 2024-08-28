package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/neo4j/cli/common/clictx"
)

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

	url, err := config.GetString("aura.auth-url")

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
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

	credential.UpdateAccessToken(grant.AccessToken, grant.ExpiresIn)
	config.Write()
	return grant.AccessToken, nil
}
