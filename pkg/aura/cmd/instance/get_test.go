package instance_test

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

func TestGetInstance(t *testing.T) {
	assert := assert.New(t)

	var instanceId = "2f49c2b3"

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":"3600","token_type":"bearer"}`))
	})

	var getCounter = 0
	mux.HandleFunc(fmt.Sprintf("/v1/instances/%s", instanceId), func(res http.ResponseWriter, req *http.Request) {
		getCounter++

		assert.Equal(http.MethodGet, req.Method)
		assert.Equal(fmt.Sprintf("/v1/instances/%s", instanceId), req.URL.Path)

		res.WriteHeader(200)
		res.Write([]byte(`{
			"data": {
				"id": "2f49c2b3",
				"name": "Production",
				"status": "running",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"connection_url": "YOUR_CONNECTION_URL",
				"metrics_integration_url": "YOUR_METRICS_INTEGRATION_ENDPOINT",
				"region": "europe-west1",
				"type": "enterprise-db",
				"memory": "8GB",
				"storage": "16GB"
			}
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"instance", "get", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), instanceId})

	cfg, err := clicfg.NewConfigFrom(strings.NewReader(`{}`), nil)
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
		"id": "2f49c2b3",
		"name": "Production",
		"status": "running",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"connection_url": "YOUR_CONNECTION_URL",
		"metrics_integration_url": "YOUR_METRICS_INTEGRATION_ENDPOINT",
		"region": "europe-west1",
		"type": "enterprise-db",
		"memory": "8GB",
		"storage": "16GB"
	}
}
`, string(out))
}
