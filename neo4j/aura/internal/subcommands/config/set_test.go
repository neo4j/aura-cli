package config_test

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

func TestSetConfig(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.NewCmd("test")
	cmd.SetArgs([]string{"config", "set", "auth-url", "test"})

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

	assert.Equal(fmt.Sprintf(`{"aura":{"base-url":"%s","auth-url":"test","output":"default","credentials":[]}}`, clicfg.DefaultAuraBaseUrl), out)
}
