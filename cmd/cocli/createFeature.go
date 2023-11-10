package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var createFeatureAppId int
var createFeatureCommitSha string
var createFeatureName string

var createFeatureCmd = &cobra.Command{
	Use:   "create <branch_name>",
	Short: "Create a new coherence feature",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		if createFeatureName == "" {
			createFeatureName = branchName
		}

		createFeaturePayloadData := map[string]string{
			"application_id": fmt.Sprint(createFeatureAppId),
			"branch_name":    branchName,
			"commit_sha":     createFeatureCommitSha,
			"name":           createFeatureName,
		}
		payloadBytes, err := json.Marshal(createFeaturePayloadData)
		if err != nil {
			panic(err)
		}
		payload := bytes.NewBuffer(payloadBytes)

		featuresCreateUrl := fmt.Sprintf(
			"https://%s%s/features",
			cocli.GetCoherenceDomain(),
			cocli.GetCoherenceApiPrefix(),
		)
		res, err := cocli.CoherenceApiRequest(
			"POST",
			featuresCreateUrl,
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
	createFeatureCmd.Flags().IntVarP(&createFeatureAppId, "app_id", "a", 0, "App ID (required)")
	createFeatureCmd.MarkFlagRequired("app_id")
	createFeatureCmd.Flags().StringVarP(&createFeatureCommitSha, "commit_sha", "c", "", "Commit SHA (optional)")
	createFeatureCmd.Flags().StringVarP(&createFeatureName, "name", "n", "", "Feature name (optional - defaults to branch_name)")
}
