package token

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/api"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewCreateCmd(cfg *clicfg.Config) *cobra.Command {
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
		Use:   "create",
		Short: "Register a deployment",
		Long:  "This endpoint registers a Fleet Manager deployment and returns a token that can be used to activate Fleet Manager in a Neo4j database.",
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
				Method:   http.MethodPost,
				PostBody: body,
				Version:  api.AuraApiVersion2,
			})
			if err != nil {
				return err
			}

			// NOTE: Deployment register should not return OK (200), it always returns 201, checking both just in case
			if statusCode == http.StatusCreated || statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"token"})
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&organizationId, organizationIdFlag, "o", "", "(required) Organization ID")
	cmd.Flags().StringVarP(&projectId, projectIdFlag, "p", "", "(required) Project/tenant ID")
	cmd.Flags().StringVarP(&deploymentId, deploymentIdFlag, "d", "", "(required) Deployment ID")
	cmd.Flags().BoolVarP(&noAutoRotate, noAutoRotateFlag, "r", false, "An optional argument to prevent the token from auto rotating when it expires.")
	cmd.Flags().StringVarP(&expiresIn, expiresInFlag, "e", "", "An optional expires in time. Accepted values are '15 minutes', '3 months', '6 months', '9 months' and '12 months'")

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

func validateExpiresIn(expiresIn string) error {
	if expiresIn != "15 minutes" && expiresIn != "3 months" && expiresIn != "6 months" && expiresIn != "9 months" && expiresIn != "12 months" {
		return fmt.Errorf("incorrect expires-in value, must be one of '15 minutes', '3 months', '6 months', '9 months' and '12 months'")
	}
	return nil
}
