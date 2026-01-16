package setting_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestRemoveSetting(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetSettingsValue("aura.settings", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})

	helper.ExecuteCommand("setting remove test")

	helper.AssertSettingsValue("aura.settings", "[]")
}
