// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clicfg/credentials"
	"github.com/neo4j/cli/common/redact"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestUnrecognisedStatusDoesNotLeakSecretIntoPanicMessage(t *testing.T) {
	rawArgs := []string{"credential", "add", "--name", "demo", "--client-id", "abc", "--client-secret", "S3cr3t-Value"}
	redactedArgs := redact.Args(rawArgs)
	redact.SetCapturedArgs(redactedArgs)
	defer redact.SetCapturedArgs(nil)

	res := &http.Response{
		StatusCode: http.StatusUnprocessableEntity, // 422 — not one of the explicitly handled codes
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"errors":[{"message":"bad request"}]}`))),
	}

	cfg := clicfg.NewConfig(afero.NewMemMapFs(), "test")
	credential := &credentials.AuraCredential{Name: "demo"}

	var panicMessage string
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicMessage = fmt.Sprint(r)
			}
		}()
		_ = handleResponseError(res, credential, cfg)
	}()

	assert.NotEmpty(t, panicMessage, "expected handleResponseError to panic on an unrecognised status code")
	assert.NotContains(t, panicMessage, "S3cr3t-Value",
		"the client secret from the command line must never appear in a panic/error message that tells the user to file a public GitHub issue")
	assert.Contains(t, panicMessage, "****",
		"the client secret should be masked with **** in the panic message")
	assert.Contains(t, panicMessage, "demo",
		"safe flag values like 'demo' should appear unmasked in the panic message")
	assert.Contains(t, panicMessage, "please report an issue")
}
