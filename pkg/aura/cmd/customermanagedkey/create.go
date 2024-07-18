package customermanagedkey

import (
	"github.com/neo4j/cli/pkg/aura/api"
	"github.com/spf13/cobra"
)

var (
	region        string
	name          string
	instanceType  string
	tenantId      string
	cloudProvider string
	keyId         string
)

const (
	regionFlag        = "region"
	nameFlag          = "name"
	instanceTypeFlag  = "type"
	tenantIdFlag      = "tenant-id"
	cloudProviderFlag = "cloud-provider"
	keyIdFlag         = "key-id"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new customer managed key",
	Long: `This subcommand creates a new Customer Managed Key in Aura. Creating a new key is an asynchronous operation.

Before you can use the key you will need to setup permissions for it. Log in to the Console, navigate to 'Customer Managed Keys' and click on the Edit icon next to the Key in order to see the instructions.

You can poll the current status of this operation by periodically getting the key details using the get subcommand.

Once the key has a status of ready you can use it for creating new instances by setting the --customer-managed-key-id flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		body := map[string]any{
			"region":         region,
			"name":           name,
			"instance_type":  instanceType,
			"tenant_id":      tenantId,
			"cloud_provider": cloudProvider,
			"key_id":         keyId,
		}

		return api.MakeRequest(cmd, "POST", "/customer-managed-keys", body)
	},
}

func init() {
	CreateCmd.Flags().StringVar(&region, regionFlag, "", "The region where the instance is hosted.")
	CreateCmd.MarkFlagRequired(regionFlag)

	CreateCmd.Flags().StringVar(&name, nameFlag, "", "The name of the instance (any UTF-8 characters with no trailing or leading whitespace).")
	CreateCmd.MarkFlagRequired(nameFlag)

	CreateCmd.Flags().StringVar(&instanceType, instanceTypeFlag, "", "The type of the instance.")
	CreateCmd.MarkFlagRequired(instanceTypeFlag)

	CreateCmd.Flags().StringVar(&tenantId, tenantIdFlag, "", "The tenant??????????")

	CreateCmd.Flags().StringVar(&cloudProvider, cloudProviderFlag, "", "The cloud provider hosting the instance.")
	CreateCmd.MarkFlagRequired(cloudProviderFlag)

	CreateCmd.Flags().StringVar(&keyId, keyIdFlag, "", "The cloud provider hosting the instance.")
	CreateCmd.MarkFlagRequired(keyIdFlag)
}
