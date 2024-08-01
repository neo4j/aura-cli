package instance_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neo4j/cli/internal/testutils"
	"github.com/neo4j/cli/pkg/aura"
	"github.com/neo4j/cli/pkg/clicfg"
	"github.com/neo4j/cli/pkg/clictx"
	"github.com/stretchr/testify/assert"
)

func TestListInstances(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":3600,"token_type":"bearer"}`))
	})

	var listCounter = 0
	mux.HandleFunc("/v1/instances", func(res http.ResponseWriter, req *http.Request) {
		listCounter++

		assert.Equal(http.MethodGet, req.Method)
		assert.Equal("/v1/instances", req.URL.Path)

		res.WriteHeader(200)
		res.Write([]byte(`{
			"data": [
				{
					"id": "2f49c2b3",
					"name": "Production",
					"tenant_id": "YOUR_TENANT_ID",
					"cloud_provider": "gcp"
				},
				{
					"id": "b51dc964",
					"name": "Instance01",
					"tenant_id": "YOUR_TENANT_ID",
					"cloud_provider": "aws"
				},
				{
					"id": "432392ae",
					"name": "Recommendations",
					"tenant_id": "YOUR_TENANT_ID",
					"cloud_provider": "azure"
				},
				{
					"id": "524b7d8d",
					"name": "Northwind",
					"tenant_id": "YOUR_TENANT_ID",
					"cloud_provider": "gcp"
				}
			]
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"instance", "list", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL)})

	fs, err := testutils.GetTestFs(`{
		"aura": {
			"credentials": [{
				"name": "test-cred",
				"access-token": "dsa",
				"token-expiry": 123
			}],
			"default-credential": "test-cred"
		}
	}`)
	assert.Nil(err)

	cfg, err := clicfg.NewConfig(fs)
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
			"id": "2f49c2b3",
			"name": "Production",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp"
		},
		{
			"id": "b51dc964",
			"name": "Instance01",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "aws"
		},
		{
			"id": "432392ae",
			"name": "Recommendations",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "azure"
		},
		{
			"id": "524b7d8d",
			"name": "Northwind",
			"tenant_id": "YOUR_TENANT_ID",
			"cloud_provider": "gcp"
		}
	]
}
`, string(out))
}
