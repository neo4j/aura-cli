package deployment

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewRefreshCmd(cfg *clicfg.Config) *cobra.Command {
	var (
		organizationId string
		projectId      string
		deploymentId   string
		expiresIn      string
		noAutoRotate   bool
	)

	const (
		organizationIdFlag = "organization-id"
		projectIdFlag      = "project-id"
		deploymentIdFlag   = "deployment-id"
		expiresInFlag      = "expires-in"
		noAutoRotateFlag   = "no-auto-rotate"
	)

	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh the deployment token",
		Long:  "This endpoint refreshes a Fleet Manager deployment token.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/organizations/%s/projects/%s/deployments/%s/token", organizationId, projectId, deploymentId)

			body := map[string]any{}
			if expiresIn != "" {
				err := validateExpiresIn(expiresIn)
				if err != nil {
					log.Fatal(err)
				}
				body["expires_in"] = expiresIn
			}
			if noAutoRotate {
				body["auto_rotate"] = false
			}

			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:   http.MethodPatch,
				PostBody: body,
				Version:  api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}
			spew.Dump(resBody)

			// NOTE: Token refresh should not return OK (200), it always returns 202, checking both just in case
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"token"})
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&organizationId, organizationIdFlag, "", "(required) Organization ID")
	cmd.Flags().StringVar(&projectId, projectIdFlag, "", "(required) Project/tenant ID")
	cmd.Flags().StringVar(&deploymentId, deploymentIdFlag, "", "(required) Deployment ID")
	cmd.Flags().BoolVar(&noAutoRotate, noAutoRotateFlag, false, "An optional argument to prevent the token from auto rotating when it expires.")
	cmd.Flags().StringVar(&expiresIn, expiresInFlag, "", "An optional expires in time. Accepted values are '15 minutes', '3 months', '6 months', '9 months' and '12 months'")

	err := cmd.MarkFlagRequired(organizationIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(projectIdFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired(deploymentIdFlag)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}
