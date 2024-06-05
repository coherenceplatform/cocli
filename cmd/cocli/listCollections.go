package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var listCollectionsAppId int

var listCollectionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List environment collections",
	Long:  "List all environment collections for the specified application.",
	Run: func(cmd *cobra.Command, args []string) {
		collectionsListUrl := fmt.Sprintf(
			"%s/api/v1/applications/%s/collections",
			cocli.GetCoherenceDomain(),
			fmt.Sprint(listCollectionsAppId),
		)
		res, err := cocli.CoherenceApiRequest(
			"GET",
			collectionsListUrl,
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
	listCollectionsCmd.Flags().IntVarP(&listCollectionsAppId, "app_id", "a", 0, "App ID (required)")
	listCollectionsCmd.MarkFlagRequired("app_id")
}
