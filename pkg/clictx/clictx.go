package clictx

import (
	"context"

	"github.com/neo4j/cli/pkg/clicfg"
)

type key string

var configKey = key("config")
var versionKey = key("version")

func NewContext(ctx context.Context, config *clicfg.Config, version string) (context.Context, error) {
	ctx = context.WithValue(ctx, configKey, config)
	ctx = context.WithValue(ctx, versionKey, version)

	return ctx, nil
}

func Config(ctx context.Context) (*clicfg.Config, bool) {
	config, ok := ctx.Value(configKey).(*clicfg.Config)
	return config, ok
}

func Version(ctx context.Context) (string, bool) {
	version, ok := ctx.Value(versionKey).(string)
	return version, ok
}
