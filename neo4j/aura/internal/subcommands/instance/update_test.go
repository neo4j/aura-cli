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
	"github.com/neo4j/cli/test/utils/testfs"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMemory(t *testing.T) {
	assert := assert.New(t)

	var instanceId = "2f49c2b3"

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":3600,"token_type":"bearer"}`))
	})

	var patchCounter = 0
	mux.HandleFunc(fmt.Sprintf("/v1/instances/%s", instanceId), func(res http.ResponseWriter, req *http.Request) {
		patchCounter++

		assert.Equal(http.MethodPatch, req.Method)
		assert.Equal(fmt.Sprintf("/v1/instances/%s", instanceId), req.URL.Path)
		body, err := io.ReadAll(req.Body)
		assert.Nil(err)
		assert.Equal(`{"memory":"8GB"}`, string(body))

		res.WriteHeader(200)
		res.Write([]byte(`{
	"data": {
    	"id": "2f49c2b3",
		"name": "Production",
		"status": "updating",
		"connection_url": "YOUR_CONNECTION_URL",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"memory": "8GB",
		"region": "europe-west1",
		"type": "enterprise-db"
  	}
}`))

	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"instance", "update", instanceId, "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), "--memory", "8GB"})

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
	assert.Equal(1, patchCounter)

	assert.Equal(`{
	"data": {
		"id": "2f49c2b3",
		"name": "Production",
		"status": "updating",
		"connection_url": "YOUR_CONNECTION_URL",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"memory": "8GB",
		"region": "europe-west1",
		"type": "enterprise-db"
	}
}
`, string(out))
}

func TestUpdateName(t *testing.T) {
	assert := assert.New(t)

	var instanceId = "2f49c2b3"

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":3600,"token_type":"bearer"}`))
	})

	var patchCounter = 0
	mux.HandleFunc(fmt.Sprintf("/v1/instances/%s", instanceId), func(res http.ResponseWriter, req *http.Request) {
		patchCounter++

		assert.Equal(http.MethodPatch, req.Method)
		assert.Equal(fmt.Sprintf("/v1/instances/%s", instanceId), req.URL.Path)
		body, err := io.ReadAll(req.Body)
		assert.Nil(err)
		assert.Equal(`{"name":"New Name"}`, string(body))

		res.WriteHeader(200)
		res.Write([]byte(`{
	"data": {
    	"id": "2f49c2b3",
		"name": "New Name",
		"status": "updating",
		"connection_url": "YOUR_CONNECTION_URL",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"memory": "4GB",
		"region": "europe-west1",
		"type": "enterprise-db"
  	}
}`))

	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"instance", "update", instanceId, "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), "--name", "New Name"})

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
	assert.Equal(1, patchCounter)

	assert.Equal(`{
	"data": {
		"id": "2f49c2b3",
		"name": "New Name",
		"status": "updating",
		"connection_url": "YOUR_CONNECTION_URL",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"memory": "4GB",
		"region": "europe-west1",
		"type": "enterprise-db"
	}
}
`, string(out))
}

func TestUpdateMemoryAndName(t *testing.T) {
	assert := assert.New(t)

	var instanceId = "2f49c2b3"

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"12345678","expires_in":3600,"token_type":"bearer"}`))
	})

	var patchCounter = 0
	mux.HandleFunc(fmt.Sprintf("/v1/instances/%s", instanceId), func(res http.ResponseWriter, req *http.Request) {
		patchCounter++

		assert.Equal(http.MethodPatch, req.Method)
		assert.Equal(fmt.Sprintf("/v1/instances/%s", instanceId), req.URL.Path)
		body, err := io.ReadAll(req.Body)
		assert.Nil(err)
		assert.Equal(`{"memory":"8GB","name":"New Name"}`, string(body))

		res.WriteHeader(200)
		res.Write([]byte(`{
	"data": {
    	"id": "2f49c2b3",
		"name": "New Name",
		"status": "updating",
		"connection_url": "YOUR_CONNECTION_URL",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"memory": "8GB",
		"region": "europe-west1",
		"type": "enterprise-db"
  	}
}`))

	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"instance", "update", instanceId, "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL), "--name", "New Name", "--memory", "8GB"})

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
	assert.Equal(1, patchCounter)

	assert.Equal(`{
	"data": {
		"id": "2f49c2b3",
		"name": "New Name",
		"status": "updating",
		"connection_url": "YOUR_CONNECTION_URL",
		"tenant_id": "YOUR_TENANT_ID",
		"cloud_provider": "gcp",
		"memory": "8GB",
		"region": "europe-west1",
		"type": "enterprise-db"
	}
}
`, string(out))
}

func TestUpdateErrorsWithNoFlags(t *testing.T) {
	assert := assert.New(t)

	var instanceId = "2f49c2b3"

	mux := http.NewServeMux()

	var authCounter = 0
	mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		authCounter++
	})

	var patchCounter = 0
	mux.HandleFunc(fmt.Sprintf("/v1/instances/%s", instanceId), func(res http.ResponseWriter, req *http.Request) {
		patchCounter++
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.SetArgs([]string{"instance", "update", instanceId, "--auth-url", fmt.Sprintf("%s/oauth/token", server.URL), "--base-url", fmt.Sprintf("%s/v1", server.URL)})

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
	assert.ErrorContains(err, `at least one of the flags in the group [memory name] is required`)

	assert.Equal(0, authCounter)
	assert.Equal(0, patchCounter)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(`Error: at least one of the flags in the group [memory name] is required
Usage:
  aura instance update [flags]

Flags:
  -h, --help            help for update
      --memory string   The size of the instance memory in GB.
      --name string     The name of the instance (any UTF-8 characters with no trailing or leading whitespace).

Global Flags:
      --auth-url string   
      --base-url string   
      --output string

`, string(out))
}
