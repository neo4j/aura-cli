package testutils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/shlex"
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
	Body   map[string]interface{}
}

type requestHandlerMock struct {
	Calls []call
	t     *testing.T
}

func (mock *requestHandlerMock) AssertCalledTimes(times int) {
	calls := len(mock.Calls)

	assert.Equal(mock.t, times, calls, "Request handler mock not called the expected number of times")
}

func (mock *requestHandlerMock) AssertCalledWithMethod(method string) {
	methods := ""

	for _, call := range mock.Calls {
		if call.Method == method {
			return
		}

		methods += call.Method
	}

	assert.Fail(mock.t, fmt.Sprintf("Handler not called with method:\nexpected: %s, actual: %s", method, methods))
}

func (mock *requestHandlerMock) AssertCalledWithBody(body string) {
	unmarshalled, err := UmarshalJson([]byte(body))
	assert.Nil(mock.t, err)

	bodies := ""

	for _, call := range mock.Calls {
		if cmp.Equal(call.Body, unmarshalled) {
			return
		}
		data, err := MarshalJson(call.Body)
		assert.Nil(mock.t, err)

		bodies += data + "\n"
	}

	assert.Fail(mock.t, fmt.Sprintf("Handler not called with body:\nexpected: %s\nactual: %s", body, bodies))
}

type AuraTestHelper struct {
	mux    *http.ServeMux
	Server *httptest.Server
	cmd    *cobra.Command
	out    *bytes.Buffer
	err    *bytes.Buffer
	cfg    string
	t      *testing.T
}

func (helper *AuraTestHelper) AddRequestHandler(path string, handler func(res http.ResponseWriter, req *http.Request)) {
	helper.mux.HandleFunc(path, handler)
}

func (helper *AuraTestHelper) Close() {
	helper.Server.Close()
}

func (helper *AuraTestHelper) ExecuteCommand(command string) {
	args, err := shlex.Split(command)
	assert.Nil(helper.t, err)

	args = append(args, "--auth-url", fmt.Sprintf("%s/oauth/token", helper.Server.URL), "--base-url", fmt.Sprintf("%s/v1", helper.Server.URL))

	helper.cmd.SetArgs(args)

	fs, err := testfs.GetTestFs(helper.cfg)
	assert.Nil(helper.t, err)

	cfg, err := clicfg.NewConfig(fs)
	assert.Nil(helper.t, err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(helper.t, err)

	helper.cmd.ExecuteContext(ctx)
}

func (helper *AuraTestHelper) SetConfig(cfg string) {
	helper.cfg = cfg
}

func (helper *AuraTestHelper) AssertErr(expected string) {
	out, err := io.ReadAll(helper.err)
	assert.Nil(helper.t, err)

	assert.Equal(helper.t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}

func (helper *AuraTestHelper) AssertOut(expected string) {
	out, err := io.ReadAll(helper.out)
	assert.Nil(helper.t, err)

	assert.Equal(helper.t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}

func (helper *AuraTestHelper) AssertOutJson(expected string) {
	out, err := io.ReadAll(helper.out)
	assert.Nil(helper.t, err)

	formatted, err := FormatJson(expected)
	assert.Nil(helper.t, err)

	assert.Equal(helper.t, formatted, string(out))
}

func (helper *AuraTestHelper) NewRequestHandlerMock(path string, status int, body string) *requestHandlerMock {
	mock := requestHandlerMock{Calls: []call{}, t: helper.t}

	helper.mux.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
		requestBody, err := io.ReadAll(req.Body)
		assert.Nil(helper.t, err)

		var unmarshalledBody map[string]interface{}
		if len(requestBody) > 0 {
			unmarshalledBody, err = UmarshalJson(requestBody)
			assert.Nil(helper.t, err)
		}

		mock.Calls = append(mock.Calls, call{Method: req.Method, Path: req.URL.Path, Body: unmarshalledBody})

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
