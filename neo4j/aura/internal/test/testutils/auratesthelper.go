package testutils

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
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type call struct {
	Method string
	Path   string
	Body   string
}

type requestHandlerMock struct {
	Calls []call
}

type AuraTestHelper struct {
	mux    *http.ServeMux
	Server *httptest.Server
	cmd    *cobra.Command
	out    *bytes.Buffer
	err    *bytes.Buffer
	ctx    context.Context
	cfg    string
	t      *testing.T
}

func (helper *AuraTestHelper) AddRequestHandler(path string, handler func(res http.ResponseWriter, req *http.Request)) {
	helper.mux.HandleFunc(path, handler)
}

func (helper *AuraTestHelper) Close() {
	helper.Server.Close()
}

func (helper *AuraTestHelper) ExecuteCommand(args []string) {
	args = append(args, "--auth-url", fmt.Sprintf("%s/oauth/token", helper.Server.URL), "--base-url", fmt.Sprintf("%s/v1", helper.Server.URL))

	helper.cmd.SetArgs(args)

	fs, err := testfs.GetTestFs(helper.cfg)
	assert.Nil(helper.t, err)

	cfg, err := clicfg.NewConfig(fs)
	assert.Nil(helper.t, err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(helper.t, err)

	err = helper.cmd.ExecuteContext(ctx)
	assert.Nil(helper.t, err)
}

func (helper *AuraTestHelper) SetConfig(cfg string) {
	helper.cfg = cfg
}

func (helper *AuraTestHelper) AssertOut(expected string) {
	out, err := io.ReadAll(helper.out)
	assert.Nil(helper.t, err)

	assert.Equal(helper.t, expected, string(out))
}

func (helper *AuraTestHelper) NewRequestHandlerMock(path string, status int, body string) *requestHandlerMock {
	mock := requestHandlerMock{}
	mock.Calls = []call{}

	helper.mux.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
		print("BONJOUR")

		requestBody, err := io.ReadAll(req.Body)
		assert.Nil(helper.t, err)

		mock.Calls = append(mock.Calls, call{Method: req.Method, Path: req.URL.Path, Body: string(requestBody)})

		res.WriteHeader(status)
		res.Write([]byte(body))
	})

	return &mock
}

func NewAuraTestHelper(t *testing.T) AuraTestHelper {
	helper := AuraTestHelper{}

	helper.t = t

	cmd := aura.NewCmd()

	out := bytes.NewBufferString("")
	err := bytes.NewBufferString("")

	cmd.SetOut(out)
	cmd.SetErr(err)

	helper.cmd = cmd
	helper.out = out
	helper.err = err

	helper.cfg = `{
				"aura": {
			"credentials": [{
				"name": "test-cred",
				"access-token": "dsa",
				"token-expiry": 123
			}],
			"default-credential": "test-cred"
		}
	}`

	helper.mux = http.NewServeMux()
	helper.mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"<token>","expires_in":3600,"token_type":"bearer"}`))
	})

	helper.Server = httptest.NewServer(helper.mux)

	return helper
}
