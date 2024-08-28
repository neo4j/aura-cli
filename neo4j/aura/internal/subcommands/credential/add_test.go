package credential_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clictx"
	"github.com/neo4j/cli/neo4j/aura"
	"github.com/neo4j/cli/test/utils/testfs"
	"github.com/stretchr/testify/assert"
)

func TestAddCredential(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.NewCmd("test")
	cmd.SetArgs([]string{"credential", "add", "--name", "test", "--client-id", "testclientid", "--client-secret", "testclientsecret"})

	fs, err := testfs.GetDefaultTestFs()
	assert.Nil(err)

	cfg, err := clicfg.NewConfig(fs)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.Nil(err)

	out, err := testfs.GetTestConfig(fs)
	assert.Nil(err)

	assert.Equal(fmt.Sprintf(`{"aura":{"base-url":"%s","auth-url":"%s","output":"json","credentials":[{"name":"test","client-id":"testclientid","client-secret":"testclientsecret","access-token":"","token-expiry":0}]}}`, clicfg.DefaultAuraBaseUrl, clicfg.DefaultAuraAuthUrl), out)
}

func TestAddCredentialIfAlreadyExists(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.NewCmd("test")
	cmd.SetArgs([]string{"credential", "add", "--name", "test", "--client-id", "testclientid", "--client-secret", "testclientsecret"})

	fs, err := testfs.GetTestFs(`{"aura":{"credentials":[{"name":"test","client-id":"testclientid","client-secret":"testclientsecret"}]}}`)
	assert.Nil(err)

	cfg, err := clicfg.NewConfig(fs)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.ErrorContains(err, "already have credential with name test")
}
