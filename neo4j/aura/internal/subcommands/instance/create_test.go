package instance_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clictx"
	"github.com/neo4j/cli/neo4j/aura"
	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
	"github.com/neo4j/cli/test/utils/testfs"
	"github.com/stretchr/testify/assert"
)

func TestCreateFreeInstance(t *testing.T) {
	assert := assert.New(t)

	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	var postCounter = 0
	helper.AddRequestHandler("/v1/instances", func(res http.ResponseWriter, req *http.Request) {
		postCounter++

		assert.Equal(http.MethodPost, req.Method)
		assert.Equal("/v1/instances", req.URL.Path)
		body, err := io.ReadAll(req.Body)
		assert.Nil(err)
		assert.Equal(`{"cloud_provider":"gcp","memory":"1GB","name":"Instance01","region":"europe-west1","tenant_id":"YOUR_TENANT_ID","type":"free-db","version":"5"}`, string(body))

		res.WriteHeader(200)
		res.Write([]byte(`{
			"data": {
				"id": "db1d1234",
				"connection_url": "YOUR_CONNECTION_URL",
				"username": "neo4j",
				"password": "letMeIn123!",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"region": "europe-west1",
				"type": "free-db",
				"name": "Instance01"
			}
		}`))

	})

	helper.ExecuteCommand([]string{"instance", "create", "--region", "europe-west1", "--name", "Instance01", "--type", "free-db", "--tenant-id", "YOUR_TENANT_ID", "--cloud-provider", "gcp"})

	assert.Equal(1, postCounter)

	helper.AssertOut(`{
	"data": {
		"id": "db1d1234",
		"connection_url": "YOUR_CONNECTION_URL",
		"username": "neo4j",
		"password": "letMeIn123!",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"region": "europe-west1",
		"type": "free-db",
		"name": "Instance01"
	}
}
`)
}

func TestCreateProfessionalInstance(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":3600,"token_type":"bearer"}`))
	})

	var postCounter = 0
	mux.HandleFunc("/v1/instances", func(res http.ResponseWriter, req *http.Request) {
		postCounter++

		assert.Equal(http.MethodPost, req.Method)
		assert.Equal("/v1/instances", req.URL.Path)
		body, err := io.ReadAll(req.Body)
		assert.Nil(err)
		assert.Equal(`{"cloud_provider":"gcp","memory":"4GB","name":"Instance01","region":"europe-west1","tenant_id":"YOUR_TENANT_ID","type":"professional-db","version":"5"}`, string(body))

		res.WriteHeader(200)
		res.Write([]byte(`{
			"data": {
				"id": "db1d1234",
				"connection_url": "YOUR_CONNECTION_URL",
				"username": "neo4j",
				"password": "letMeIn123!",
				"tenant_id": "YOUR_TENANT_ID",
				"cloud_provider": "gcp",
				"region": "europe-west1",
				"type": "professional-db",
				"name": "Instance01"
			}
		}`))

	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"instance", "create", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), "--region", "europe-west1", "--name", "Instance01", "--type", "professional-db", "--tenant-id", "YOUR_TENANT_ID", "--cloud-provider", "gcp", "--memory", "4GB"})

	fs, err := testfs.GetTestFs(`{
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
	assert.Equal(1, postCounter)

	assert.Equal(`{
	"data": {
		"id": "db1d1234",
		"connection_url": "YOUR_CONNECTION_URL",
		"username": "neo4j",
		"password": "letMeIn123!",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"region": "europe-west1",
		"type": "professional-db",
		"name": "Instance01"
	}
}
`, string(out))
}

func TestCreateProfessionalInstanceNoMemory(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
	})

	var postCounter = 0
	mux.HandleFunc("/v1/instances", func(res http.ResponseWriter, req *http.Request) {
		postCounter++
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.SetArgs([]string{"instance", "create", "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), "--region", "europe-west1", "--name", "Instance01", "--type", "professional-db", "--tenant-id", "YOUR_TENANT_ID", "--cloud-provider", "gcp"})

	fs, err := testfs.GetTestFs(`{
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
	assert.ErrorContains(err, `required flag(s) "memory" not set`)

	assert.Equal(0, authCounter)
	assert.Equal(0, postCounter)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(`Error: required flag(s) "memory" not set
Usage:
  aura instance create [flags]

Flags:
      --cloud-provider string            The cloud provider hosting the instance.
      --customer-managed-key-id string   
  -h, --help                             help for create
      --memory string                    The size of the instance memory in GB.
      --name string                      The name of the instance (any UTF-8 characters with no trailing or leading whitespace).
      --region string                    The region where the instance is hosted.
      --tenant-id string                 
      --type string                      The type of the instance.
      --version string                   The Neo4j version of the instance. (default "5")

Global Flags:
      --auth-url string   
      --base-url string   
      --output string

`, string(out))
}
