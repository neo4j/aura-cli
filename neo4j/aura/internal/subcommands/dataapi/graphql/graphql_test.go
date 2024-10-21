package graphql_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j/aura/internal/test/testutils"
)

func TestGraphQLDataApiNotEnabledError(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.ExecuteCommand("data-api graphql list")

	helper.AssertErr("Error: The command 'data-api' is beta functionality. Turn it on by setting the Aura config key 'beta-enabled' to 'true'.")
}
