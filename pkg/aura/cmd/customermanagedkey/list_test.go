package customermanagedkey_test

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

func TestListCustomerManagedKeys(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++

		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":"3600","token_type":"bearer"}`))
	})

	var listCounter = 0
	mux.HandleFunc("/v1/customer-managed-keys", func(res http.ResponseWriter, req *http.Request) {
		listCounter++

		assert.Equal(http.MethodGet, req.Method)
		assert.Equal("/v1/customer-managed-keys", req.URL.Path)

		res.WriteHeader(200)
		res.Write([]byte(`{
		"data": [
			{
				"id": "f15cc45b-1c29-44e8-911f-3ba719f70ed7",
				"name": "Production Key",
				"tenant_id": "YOUR_TENANT_ID"
			},
			{
				"id": "0d971cc4-f703-40fd-8c5c-f5ec134f6c84",
				"name": "Dev Key",
				"tenant_id": "YOUR_TENANT_ID"
			}
		]
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"customer-managed-key", "list", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL)})

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
			"id": "f15cc45b-1c29-44e8-911f-3ba719f70ed7",
			"name": "Production Key",
			"tenant_id": "YOUR_TENANT_ID"
		},
		{
			"id": "0d971cc4-f703-40fd-8c5c-f5ec134f6c84",
			"name": "Dev Key",
			"tenant_id": "YOUR_TENANT_ID"
		}
	]
}
`, string(out))
}

func TestListCustomerManagedKeysAlias(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++

		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":"3600","token_type":"bearer"}`))
	})

	var listCounter = 0
	mux.HandleFunc("/v1/customer-managed-keys", func(res http.ResponseWriter, req *http.Request) {
		listCounter++

		assert.Equal(http.MethodGet, req.Method)
		assert.Equal("/v1/customer-managed-keys", req.URL.Path)

		res.WriteHeader(200)
		res.Write([]byte(`{
		"data": [
			{
				"id": "f15cc45b-1c29-44e8-911f-3ba719f70ed7",
				"name": "Production Key",
				"tenant_id": "YOUR_TENANT_ID"
			},
			{
				"id": "0d971cc4-f703-40fd-8c5c-f5ec134f6c84",
				"name": "Dev Key",
				"tenant_id": "YOUR_TENANT_ID"
			}
		]
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"cmk", "list", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL)})

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
			"id": "f15cc45b-1c29-44e8-911f-3ba719f70ed7",
			"name": "Production Key",
			"tenant_id": "YOUR_TENANT_ID"
		},
		{
			"id": "0d971cc4-f703-40fd-8c5c-f5ec134f6c84",
			"name": "Dev Key",
			"tenant_id": "YOUR_TENANT_ID"
		}
	]
}
`, string(out))
}
