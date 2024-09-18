package config_test

import (
	"fmt"
	"testing"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestListConfig(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.OverwriteConfig("{}")

	helper.ExecuteCommand("config list")

	helper.AssertOutJson(fmt.Sprintf(`{"auth-url": "%s","base-url": "%s","credentials": [],"output": "json", "polling": {"interval": 30, "max-retries": 20}}`, clicfg.DefaultAuraAuthUrl, clicfg.DefaultAuraBaseUrl))
}
