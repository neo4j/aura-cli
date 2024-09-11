package testutils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/shlex"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura"
	"github.com/neo4j/cli/test/utils/testfs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type AuraTestHelper struct {
	mux    *http.ServeMux
	Server *httptest.Server
	out    *bytes.Buffer
	err    *bytes.Buffer
	cfg    string
	fs     afero.Fs
	t      *testing.T
}

func (helper *AuraTestHelper) Close() {
	helper.Server.Close()
}

func (helper *AuraTestHelper) ExecuteCommand(command string) {
	args, err := shlex.Split(command)
	assert.Nil(helper.t, err)

	fs, err := testfs.GetTestFs(helper.cfg)
	assert.Nil(helper.t, err)

	helper.fs = fs

	cfg, err := clicfg.NewConfig(fs, "test")
	assert.Nil(helper.t, err)

	cfg.Aura.SetPollingConfig(2, 0)

	cmd := aura.NewCmd(cfg)

	cmd.SetArgs(args)

	cmd.SetOut(helper.out)
	cmd.SetErr(helper.err)

	cmd.Execute()
}

func (helper *AuraTestHelper) SetConfig(cfg string) {
	helper.cfg = cfg
}

func (helper *AuraTestHelper) OverwriteConfig(cfg string) {
	helper.cfg = cfg
}

func (helper *AuraTestHelper) SetConfigValue(key string, value interface{}) {
	cfg, err := sjson.Set(helper.cfg, key, value)
	assert.Nil(helper.t, err)
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

func (helper *AuraTestHelper) PrintOut() string {
	out, err := io.ReadAll(helper.out)
	assert.Nil(helper.t, err)

	return string(out)
}
func (helper *AuraTestHelper) PrintErr() string {
	out, err := io.ReadAll(helper.err)
	assert.Nil(helper.t, err)

	return string(out)
}

func (helper *AuraTestHelper) AssertOutJson(expected string) {
	out, err := io.ReadAll(helper.out)
	assert.Nil(helper.t, err)

	formatted, err := FormatJson(expected, "\t")
	assert.Nil(helper.t, err)

	assert.Equal(helper.t, formatted, string(out))
}

func (helper *AuraTestHelper) AssertConfig(expected string) {
	file, err := helper.fs.Open(filepath.Join(clicfg.ConfigPrefix, "neo4j", "cli", "config.json"))
	assert.Nil(helper.t, err)
	defer file.Close()

	out, err := io.ReadAll(file)
	assert.Nil(helper.t, err)

	formatted, err := FormatJson(expected, "  ")
	assert.Nil(helper.t, err)

	assert.Equal(helper.t, formatted, string(out))
}

func (helper *AuraTestHelper) AssertConfigValue(key string, expected string) {
	file, err := helper.fs.Open(filepath.Join(clicfg.ConfigPrefix, "neo4j", "cli", "config.json"))
	assert.Nil(helper.t, err)
	defer file.Close()

	out, err := io.ReadAll(file)
	assert.Nil(helper.t, err)

	actual := gjson.Get(string(out), key)

	formattedExpected, err := FormatJson(expected, "\t")
	if err != nil {
		formattedExpected = expected
	}

	formattedActual, err := FormatJson(actual.String(), "\t")
	if err != nil {
		formattedActual = actual.String()
	}

	assert.Equal(helper.t, formattedExpected, formattedActual)
}

func (helper *AuraTestHelper) NewRequestHandlerMock(path string, status int, body string) *requestHandlerMock {
	mock := requestHandlerMock{Calls: []call{}, t: helper.t, Responses: []response{
		{status: status, body: body},
	}}

	helper.mux.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
		requestBody, err := io.ReadAll(req.Body)
		assert.Nil(helper.t, err)

		var unmarshalledBody map[string]interface{}
		if len(requestBody) > 0 {
			unmarshalledBody, err = UmarshalJson(requestBody)
			assert.Nil(helper.t, err)
		}

		requestCount := len(mock.Calls)
		mock.Calls = append(mock.Calls, call{Method: req.Method, Path: req.URL.Path, Body: unmarshalledBody})

		if requestCount >= len(mock.Responses) {
			res.WriteHeader(404)
		} else {
			response := mock.Responses[requestCount]

			res.WriteHeader(response.status)
			res.Write([]byte(response.body))
		}
	})

	return &mock
}

func NewAuraTestHelper(t *testing.T) AuraTestHelper {
	helper := AuraTestHelper{}

	helper.t = t

	helper.out = bytes.NewBufferString("")
	helper.err = bytes.NewBufferString("")

	helper.mux = http.NewServeMux()
	helper.mux.HandleFunc("/oauth/token", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"<token>","expires_in":3600,"token_type":"bearer"}`))
	})

	server := httptest.NewServer(helper.mux)

	helper.cfg = fmt.Sprintf(`{
				"aura": {
					"auth-url": "%s/oauth/token",
					"base-url": "%s/v1",
					"credentials": [{
						"name": "test-cred",
						"access-token": "dsa",
						"token-expiry": 123
					}],
					"default-credential": "test-cred"
					}
				}`, server.URL, server.URL)

	helper.Server = server

	return helper
}
