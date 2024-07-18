package config_test

import (
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

func TestGetConfig(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.Cmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"config", "get", "auth-url"})

	cfg, err := clicfg.NewConfigFrom(strings.NewReader(`{"aura":{"auth-url":"test"}}`), nil)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")

	assert.Nil(err)

	err = aura.Execute(ctx)
	assert.Nil(err)

	out, err := io.ReadAll(b)

	assert.Nil(err)

	assert.Equal("test\n", string(out))
}

func TestGetConfigDefault(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.Cmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"config", "get", "auth-url"})

	cfg, err := clicfg.NewConfigFrom(strings.NewReader("{}"), nil)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = aura.Execute(ctx)
	assert.Nil(err)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(fmt.Sprintf("%s\n", clicfg.DefaultAuraAuthUrl), string(out))
}
