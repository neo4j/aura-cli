package customermanagedkey

import (
	"encoding/json"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewCreateCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		region        string
		name          string
		instanceType  string
		tenantId      string
		cloudProvider string
		keyId         string
		await         bool
	)

	const (
		regionFlag        = "region"
		nameFlag          = "name"
		instanceTypeFlag  = "type"
		tenantIdFlag      = "tenant-id"
		cloudProviderFlag = "cloud-provider"
		keyIdFlag         = "key-id"
		awaitFlag         = "await"
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new customer managed key",
		Long: `This subcommand creates a new Customer Managed Key in Aura. Creating a new key is an asynchronous operation.

Before you can use the key you will need to setup permissions for it. Log in to the Console, navigate to 'Customer Managed Keys' and click on the Edit icon next to the Key in order to see the instructions.

You can poll the current status of this operation by periodically getting the key details using the get subcommand.

Once the key has a status of ready you can use it for creating new instances by setting the --customer-managed-key-id flag.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Aura.DefaultTenant() == "" {
				cmd.MarkFlagRequired(tenantIdFlag)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{
				"region":         region,
				"name":           name,
				"instance_type":  instanceType,
				"cloud_provider": cloudProvider,
				"key_id":         keyId,
			}

			if tenantId == "" {
				body["tenant_id"] = cfg.Aura.DefaultTenant()
			} else {
				body["tenant_id"] = tenantId
			}

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, http.MethodPost, "/customer-managed-keys", body)
			if err != nil {
				return err
			}
			// NOTE: Instance delete should not return OK (200), it always returns 202
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "tenant_id", "status", "created", "cloud_provider", "key_id", "region", "type"})
				if err != nil {
					return err
				}

				if await {
					cmd.Println("Waiting for customer managed key to be ready...")
					var response api.CreateCMKResponse
					if err := json.Unmarshal(resBody, &response); err != nil {
						return err
					}

					pollResponse, err := api.PollCMK(cfg, response.Data.Id, api.CMKStatusPending)
					if err != nil {
						return err
					}

					cmd.Println("CMK Status:", pollResponse.Data.Status)
				}

			}

			return nil

		},
	}

	cmd.Flags().StringVar(&region, regionFlag, "", "The region where the instance is hosted.")
	cmd.MarkFlagRequired(regionFlag)

	cmd.Flags().StringVar(&name, nameFlag, "", "The name of the instance (any UTF-8 characters with no trailing or leading whitespace).")
	cmd.MarkFlagRequired(nameFlag)

	cmd.Flags().StringVar(&instanceType, instanceTypeFlag, "", "The type of the instance.")
	cmd.MarkFlagRequired(instanceTypeFlag)

	cmd.Flags().StringVar(&tenantId, tenantIdFlag, "", "")

	cmd.Flags().StringVar(&cloudProvider, cloudProviderFlag, "", "The cloud provider hosting the instance.")
	cmd.MarkFlagRequired(cloudProviderFlag)

	cmd.Flags().StringVar(&keyId, keyIdFlag, "", "")
	cmd.MarkFlagRequired(keyIdFlag)

	cmd.Flags().BoolVar(&await, awaitFlag, false, "Waits until created customer managed key is ready.")

	return cmd
}
