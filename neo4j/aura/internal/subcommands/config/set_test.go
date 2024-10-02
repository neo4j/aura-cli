package config_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestSetConfig(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.OverwriteConfig("{}")

	helper.ExecuteCommand("config set auth-url test")

	helper.AssertConfigValue("aura.auth-url", "test")
}

func TestSetConfigWithInvalidConfigKey(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.OverwriteConfig("{}")

	helper.ExecuteCommand("config set invalid test")

	helper.AssertErr("Error: invalid config key specified: invalid")
}
