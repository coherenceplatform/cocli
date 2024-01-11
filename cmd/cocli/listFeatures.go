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
	Short: "List coherence features",
	Long:  "List all coherence features for the specified application.",
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
