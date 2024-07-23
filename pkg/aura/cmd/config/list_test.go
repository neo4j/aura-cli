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

func TestListConfigDefault(t *testing.T) {
	assert := assert.New(t)

	cmd := aura.NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"config", "list"})

	cfg, err := clicfg.NewConfigFrom(strings.NewReader("{}"), nil)
	assert.Nil(err)

	ctx, err := clictx.NewContext(context.Background(), cfg, "test")
	assert.Nil(err)

	err = cmd.ExecuteContext(ctx)
	assert.Nil(err)

	out, err := io.ReadAll(b)
	assert.Nil(err)

	assert.Equal(fmt.Sprintf("{\n\t\"base-url\": \"%s\",\n\t\"auth-url\": \"%s\",\n\t\"output\": \"json\",\n\t\"credentials\": []\n}\n", clicfg.DefaultAuraBaseUrl, clicfg.DefaultAuraAuthUrl), string(out))
}
