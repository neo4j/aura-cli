package clicfg_test

import (
	"fmt"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/test/utils/testfs"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAuraBaseUrlConfig(t *testing.T) {
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
	assert.Equal(t, server.URL, cfg.Aura.BaseUrl())
}
