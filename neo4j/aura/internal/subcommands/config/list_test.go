package config_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clictx"
	"github.com/neo4j/cli/neo4j/aura"
	"github.com/neo4j/cli/test/utils/testfs"
	"github.com/stretchr/testify/assert"
)

func TestListConfigDefault(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.NewCmd("test")
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"config", "list"})

	fs, err := testfs.GetDefaultTestFs()
	assert.Nil(err)

	cfg, err := clicfg.NewConfig(fs)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.Nil(err)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(fmt.Sprintf("{\n\t\"base-url\": \"%s\",\n\t\"auth-url\": \"%s\",\n\t\"output\": \"default\",\n\t\"credentials\": []\n}\n", clicfg.DefaultAuraBaseUrl, clicfg.DefaultAuraAuthUrl), string(out))
}
