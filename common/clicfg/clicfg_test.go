// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package clicfg_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/test/utils/testfs"
	"github.com/stretchr/testify/assert"
)

func TestGetAuraBaseUrlConfigRemovesTrailingPath(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	cfgStr := fmt.Sprintf(`{
		"aura": {
			"auth-url": "%s/oauth/token",
			"base-url": "%s/v1",
			"output": "json"
			}
		}`, server.URL, server.URL)

	credentialsStr := `{
		"aura": {
			"credentials": [{
				"name": "test-cred",
				"access-token": "dsa",
				"token-expiry": 123
			}],
			"default-credential": "test-cred"
			}
		}`

	fs, err := testfs.GetTestFs(cfgStr, credentialsStr)
	assert.Nil(t, err)
	cfg := clicfg.NewConfig(fs, "test")

	//The path parameter will be removed from GET base url
	baseUrl, err := cfg.Aura.BaseUrl()
	assert.Nil(t, err)
	assert.Equal(t, server.URL, baseUrl)
}

func TestGetAuraBaseUrlConfigReturnsErrorForInvalidUrl(t *testing.T) {
	cfgStr := `{
		"aura": {
			"auth-url": "http://example.com/oauth/token",
			"base-url": "http://exa mple.com",
			"output": "json"
			}
		}`

	credentialsStr := `{
		"aura": {
			"credentials": [{
				"name": "test-cred",
				"access-token": "dsa",
				"token-expiry": 123
			}],
			"default-credential": "test-cred"
			}
		}`

	fs, err := testfs.GetTestFs(cfgStr, credentialsStr)
	assert.Nil(t, err)
	cfg := clicfg.NewConfig(fs, "test")

	_, err = cfg.Aura.BaseUrl()
	assert.EqualError(t, err, `parse "http://exa mple.com": invalid character " " in host name`)
}
