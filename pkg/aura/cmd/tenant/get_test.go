package tenant_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/neo4j/cli/pkg/aura"
	"github.com/neo4j/cli/pkg/clicfg"
	"github.com/neo4j/cli/pkg/clictx"
	"github.com/stretchr/testify/assert"
)

func TestGetTenant(t *testing.T) {
	assert := assert.New(t)

	var tenantId = "6981ace7-efe8-4f5c-b7c5-267b5162ce91"

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":"3600","token_type":"bearer"}`))
	})

	var getCounter = 0
	mux.HandleFunc(fmt.Sprintf("/v1/tenants/%s", tenantId), func(res http.ResponseWriter, req *http.Request) {
		getCounter++

		assert.Equal(http.MethodGet, req.Method)
		assert.Equal(fmt.Sprintf("/v1/tenants/%s", tenantId), req.URL.Path)

		res.WriteHeader(200)
		res.Write([]byte(`{
			"data": {
				"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
				"name": "Production",
				"instance_configurations": []
			}
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"tenant", "get", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), tenantId})

	cfg, err := clicfg.NewConfigFrom(strings.NewReader(`{
		"aura": {
			"credentials": [{
				"name": "test-cred",
				"access-token": "dsa",
				"token-expiry": 123
			}],
			"default-credential": "test-cred"
		}
	}`), nil)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.Nil(err)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(1, authCounter)
	assert.Equal(1, getCounter)

	assert.Equal(`{
	"data": {
		"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
		"name": "Production",
		"instance_configurations": []
	}
}
`, string(out))
}
