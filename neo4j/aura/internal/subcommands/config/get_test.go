package config_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestGetConfig(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.auth-url", "test")

	helper.ExecuteCommand("config get auth-url")

	helper.AssertOut("test")
}

func TestGetConfigDefault(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.OverwriteConfig("{}")

	helper.ExecuteCommand("config get output")

	helper.AssertOut("default")
}
