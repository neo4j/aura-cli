package setting_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestUseSetting(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetSettingsValue("aura.settings", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})

	helper.ExecuteCommand("setting use test")

	helper.AssertSettingsValue("aura.default-setting", "test")
}

func TestUseSettingIfDoesNotExist(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.ExecuteCommand("setting use test")

	helper.AssertErr("Error: could not find setting with name test")
}
