package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var promoteStaticEnvId int
var promoteStaticEnvBranchName string
var promoteStaticEnvCommitSha string

var promoteStaticEnvironmentCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy to a static coherence environment",
	Run: func(cmd *cobra.Command, args []string) {
		promoteEnvPayloadData := map[string]string{
			"branch_name": promoteStaticEnvBranchName,
			"commit_sha":  promoteStaticEnvCommitSha,
		}
		payloadBytes, err := json.Marshal(promoteEnvPayloadData)
		if err != nil {
			panic(err)
		}
		payload := bytes.NewBuffer(payloadBytes)

		promoteEnvUrl := fmt.Sprintf(
			"https://%s%s/environments/%s/deploy",
			cocli.GetCoherenceDomain(),
			cocli.GetCoherenceApiPrefix(),
			fmt.Sprint(promoteStaticEnvId),
		)
		res, err := cocli.CoherenceApiRequest(
			"PUT",
			promoteEnvUrl,
			payload,
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(cocli.FormatJSONOutput(bodyBytes))
	},
}

func init() {
	promoteStaticEnvironmentCmd.Flags().IntVarP(&promoteStaticEnvId, "environment_id", "e", 0, "Target static environment ID (required)")
	promoteStaticEnvironmentCmd.MarkFlagRequired("environment_id")

	promoteStaticEnvironmentCmd.Flags().StringVarP(&promoteStaticEnvBranchName, "branch_name", "b", "", "Branch name to deploy from (required)")
	promoteStaticEnvironmentCmd.MarkFlagRequired("branch_name")

	promoteStaticEnvironmentCmd.Flags().StringVarP(&promoteStaticEnvCommitSha, "commit_sha", "c", "", "Commit sha to deploy (required)")
	promoteStaticEnvironmentCmd.MarkFlagRequired("commit_sha")
}
