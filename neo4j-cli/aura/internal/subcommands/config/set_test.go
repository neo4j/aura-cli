package config_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
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

func TestSetConfigWithInvalidOutputValue(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.OverwriteConfig("{}")

	helper.ExecuteCommand("config set output invalid")

	helper.AssertErr("Error: invalid output value specified: invalid")
}

func TestSetBetaEnabledConfig(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.OverwriteConfig("{}")

	helper.ExecuteCommand("config set beta-enabled true")

	helper.AssertConfigValue("aura.base-url", "https://api.neo4j.io/v1beta5")
	helper.AssertConfigValue("aura.beta-enabled", "true")

	helper.ExecuteCommand("config set beta-enabled false")

	helper.AssertConfigValue("aura.base-url", "https://api.neo4j.io/v1")
	helper.AssertConfigValue("aura.beta-enabled", "false")
}
