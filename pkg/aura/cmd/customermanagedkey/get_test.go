package customermanagedkey_test

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

func TestGetCustomerManagedKey(t *testing.T) {
	assert := assert.New(t)

	cmkId := "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9"

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++

		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":"3600","token_type":"bearer"}`))
	})

	var getCounter = 0
	mux.HandleFunc(fmt.Sprintf("/v1/customer-managed-keys/%s", cmkId), func(res http.ResponseWriter, req *http.Request) {
		getCounter++

		assert.Equal(http.MethodGet, req.Method)
		assert.Equal(fmt.Sprintf("/v1/customer-managed-keys/%s", cmkId), req.URL.Path)

		res.WriteHeader(200)
		res.Write([]byte(`{
			"data": {
				"id": "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9",
				"name": "Instance01",
				"created": "2024-01-31T14:06:57Z",
				"cloud_provider": "aws",
				"key_id": "arn:aws:kms:us-east-1:123456789:key/11111-a222-1212-x789-1212f1212f",
				"region": "us-east-1",
				"type": "enterprise-db",
				"tenant_id": "YOUR_TENANT_ID",
				"status": "ready"
			}
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"customer-managed-key", "get", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), cmkId})

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
	assert.Equal(1, getCounter)

	assert.Equal(`{
	"data": {
		"id": "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9",
		"name": "Instance01",
		"created": "2024-01-31T14:06:57Z",
		"cloud_provider": "aws",
		"key_id": "arn:aws:kms:us-east-1:123456789:key/11111-a222-1212-x789-1212f1212f",
		"region": "us-east-1",
		"type": "enterprise-db",
		"tenant_id": "YOUR_TENANT_ID",
		"status": "ready"
	}
}
`, string(out))
}

func TestGetCustomerManagedKeyAlias(t *testing.T) {
	assert := assert.New(t)

	cmkId := "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9"

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++

		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":"3600","token_type":"bearer"}`))
	})

	var getCounter = 0
	mux.HandleFunc(fmt.Sprintf("/v1/customer-managed-keys/%s", cmkId), func(res http.ResponseWriter, req *http.Request) {
		getCounter++

		assert.Equal(http.MethodGet, req.Method)
		assert.Equal(fmt.Sprintf("/v1/customer-managed-keys/%s", cmkId), req.URL.Path)

		res.WriteHeader(200)
		res.Write([]byte(`{
			"data": {
				"id": "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9",
				"name": "Instance01",
				"created": "2024-01-31T14:06:57Z",
				"cloud_provider": "aws",
				"key_id": "arn:aws:kms:us-east-1:123456789:key/11111-a222-1212-x789-1212f1212f",
				"region": "us-east-1",
				"type": "enterprise-db",
				"tenant_id": "YOUR_TENANT_ID",
				"status": "ready"
			}
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"cmk", "get", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), cmkId})

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
	assert.Equal(1, getCounter)

	assert.Equal(`{
	"data": {
		"id": "8c764aed-8eb3-4a1c-92f6-e4ef0c7a6ed9",
		"name": "Instance01",
		"created": "2024-01-31T14:06:57Z",
		"cloud_provider": "aws",
		"key_id": "arn:aws:kms:us-east-1:123456789:key/11111-a222-1212-x789-1212f1212f",
		"region": "us-east-1",
		"type": "enterprise-db",
		"tenant_id": "YOUR_TENANT_ID",
		"status": "ready"
	}
}
`, string(out))
}
