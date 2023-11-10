package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var listFeaturesAppId int

var listFeaturesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all coherence applications",
	Long:  "List all coherence applications that are accessible by the currently authenticated user.",
	Run: func(cmd *cobra.Command, args []string) {
		featuresListUrl := fmt.Sprintf(
			"https://%s%s/features?application_id=%s",
			cocli.GetCoherenceDomain(),
			cocli.GetCoherenceApiPrefix(),
			fmt.Sprint(listFeaturesAppId),
		)
		res, err := cocli.CoherenceApiRequest(
			"GET",
			featuresListUrl,
			nil,
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
	listFeaturesCmd.Flags().IntVarP(&listFeaturesAppId, "app_id", "a", 0, "App ID (required)")
	listFeaturesCmd.MarkFlagRequired("app_id")
}
