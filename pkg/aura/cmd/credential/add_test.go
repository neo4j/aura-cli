package credential_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/neo4j/cli/pkg/aura"
	"github.com/neo4j/cli/pkg/clicfg"
	"github.com/neo4j/cli/pkg/clictx"
	"github.com/stretchr/testify/assert"
)

func TestAddCredential(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.NewCmd()
	cmd.SetArgs([]string{"credential", "add", "--name", "test", "--client-id", "testclientid", "--client-secret", "testclientsecret"})

	b := bytes.NewBufferString("")
	cfg, err := clicfg.NewConfigFrom(strings.NewReader("{}"), bufio.NewWriter(b))
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.Nil(err)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(fmt.Sprintf(`{"aura":{"base-url":"%s","auth-url":"%s","output":"json","credentials":[{"name":"test","client-id":"testclientid","client-secret":"testclientsecret"}]}}`, clicfg.DefaultAuraBaseUrl, clicfg.DefaultAuraAuthUrl), string(out))
}

func TestAddCredentialIfAlreadyExists(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.NewCmd()
	cmd.SetArgs([]string{"credential", "add", "--name", "test", "--client-id", "testclientid", "--client-secret", "testclientsecret"})

	cfg, err := clicfg.NewConfigFrom(strings.NewReader(`{"aura":{"credentials":[{"name":"test","client-id":"testclientid","client-secret":"testclientsecret"}]}}`), nil)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.ErrorContains(err, "already have credential with name test")
}
