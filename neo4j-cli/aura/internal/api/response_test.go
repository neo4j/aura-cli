// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

// Security-audit finding 01 test: asserts that handleResponseError's branches
// do not leak command-line secrets into panic messages. Redacted args are used
// instead of raw os.Args in error output.

package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clicfg/credentials"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestF01_UnrecognisedStatusDoesNotLeakSecretIntoPanicMessage(t *testing.T) {
	origArgs := os.Args
	os.Args = []string{"aura-cli", "credential", "add", "--name", "demo", "--client-id", "abc", "--client-secret", "S3cr3t-Value"}
	defer func() { os.Args = origArgs }()

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
	assert.True(t, strings.Contains(panicMessage, "please report an issue"),
		"sanity check that this is indeed the finding-01 code path")
}
