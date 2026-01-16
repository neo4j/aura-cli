package setting_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestAddFirstSetting(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetSettingsValue("aura.settings", []map[string]string{})

	helper.ExecuteCommand("setting add --name test --organization-id testorganizationid --project-id testprojectid")

	helper.AssertSettingsValue("aura.settings", `[{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}]`)
	helper.AssertSettingsValue("aura.default-setting", "test")
}

func TestAddSettingIfAlreadyExists(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetSettingsValue("aura.settings", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})

	helper.ExecuteCommand("setting add --name test --organization-id testorganizationid --project-id testprojectid")

	helper.AssertErr("Error: already have setting with name test")
}
func TestAddAditionalSettings(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetSettingsValue("aura.settings", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})
	helper.SetSettingsValue("aura.default-setting", "test")

	helper.ExecuteCommand("setting add --name test-new --organization-id newtestorganizationid --project-id newtestprojectid")

	helper.AssertSettingsValue("aura.settings", `[{"name":"test","organization-id":"testorganizationid","project-id":"testprojectid"}, {"name":"test-new","organization-id":"newtestorganizationid","project-id":"newtestprojectid"}]`)
	helper.AssertSettingsValue("aura.default-setting", "test")
}
