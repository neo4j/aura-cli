package config_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/neo4j/cli/internal/testutils"
	"github.com/neo4j/cli/pkg/aura"
	"github.com/neo4j/cli/pkg/clicfg"
	"github.com/neo4j/cli/pkg/clictx"
	"github.com/stretchr/testify/assert"
)

func TestSetConfig(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.NewCmd()
	cmd.SetArgs([]string{"config", "set", "auth-url", "test"})

	fs, err := testutils.GetDefaultTestFs()
	assert.Nil(err)

	cfg, err := clicfg.NewConfig(fs)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.Nil(err)

	out, err := testutils.GetTestConfig(fs)
	assert.Nil(err)

	assert.Equal(fmt.Sprintf(`{"aura":{"base-url":"%s","auth-url":"test","output":"json","credentials":[]}}`, clicfg.DefaultAuraBaseUrl), out)
}
