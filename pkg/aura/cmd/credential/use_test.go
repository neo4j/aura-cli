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

func TestUseCredential(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.Cmd
	cmd.SetArgs([]string{"credential", "use", "test"})

	b := bytes.NewBufferString("")
	cfg, err := clicfg.NewConfigFrom(strings.NewReader(`{"aura":{"credentials":[{"name":"test","client-id":"testclientid","client-secret":"testclientsecret"}]}}`), bufio.NewWriter(b))
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = aura.Execute(ctx)
	assert.Nil(err)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(fmt.Sprintf(`{"aura":{"base-url":"%s","auth-url":"%s","output":"json","default-credential":"test","credentials":[{"name":"test","client-id":"testclientid","client-secret":"testclientsecret"}]}}`, clicfg.DefaultAuraBaseUrl, clicfg.DefaultAuraAuthUrl), string(out))
}

// TODO: currently fails when running with all tests - figure out what is going here
func TestUseCredentialIfDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.Cmd
	cmd.SetArgs([]string{"credential", "use", "test"})

	cfg, err := clicfg.NewConfigFrom(strings.NewReader(`{}`), nil)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = aura.Execute(ctx)
	assert.ErrorContains(err, "could not find credential with name test")
}
