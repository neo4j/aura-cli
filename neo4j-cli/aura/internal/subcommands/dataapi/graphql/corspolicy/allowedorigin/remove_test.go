package allowedorigin_test

import (
	"fmt"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestRemoveAllowedOriginFlagsValidation(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	instanceId := "2f49c2b3"
	dataApiId := "e157301d"
	allowedOrigin := "https://test.com"

	tests := map[string]struct {
		executedCommand string
		expectedError   string
	}{
		"missing all flags": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s", allowedOrigin),
			expectedError:   "Error: required flag(s) \"data-api-id\", \"instance-id\" not set",
		},
		"missing origin": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin remove --data-api-id %s --instance-id %s", dataApiId, instanceId),
			expectedError:   "Error: accepts 1 arg(s), received 0",
		},
		"missing data api id flag": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s --instance-id %s", allowedOrigin, instanceId),
			expectedError:   "Error: required flag(s) \"data-api-id\" not set",
		},
		"missing instance id flag": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin remove %s --data-api-id %s", allowedOrigin, dataApiId),
			expectedError:   "Error: required flag(s) \"instance-id\" not set",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper.ExecuteCommand(tt.executedCommand)
			helper.AssertErr(tt.expectedError)
		})
	}
}