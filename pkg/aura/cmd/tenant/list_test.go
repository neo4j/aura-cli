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

func TestListTenants(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++

		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":"3600","token_type":"bearer"}`))
	})

	var listCounter = 0
	mux.HandleFunc("/v1/tenants", func(res http.ResponseWriter, req *http.Request) {
		listCounter++

		assert.Equal(http.MethodGet, req.Method)
		assert.Equal("/v1/tenants", req.URL.Path)

		res.WriteHeader(200)
		res.Write([]byte(`{
			"data": [
				{
				"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
				"name": "Production"
				},
				{
				"id": "YOUR_TENANT_ID",
				"name": "Staging"
				},
				{
				"id": "da045ab3-3b89-4f45-8b96-528f2e47cd13",
				"name": "Development"
				}
			]
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.Cmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"tenant", "list", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL)})

	cfg, err := clicfg.NewConfigFrom(strings.NewReader(`{}`), nil)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.Nil(err)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(1, authCounter)
	assert.Equal(1, listCounter)

	assert.Equal(`{
	"data": [
		{
			"id": "6981ace7-efe8-4f5c-b7c5-267b5162ce91",
			"name": "Production"
		},
		{
			"id": "YOUR_TENANT_ID",
			"name": "Staging"
		},
		{
			"id": "da045ab3-3b89-4f45-8b96-528f2e47cd13",
			"name": "Development"
		}
	]
}
`, string(out))
}
