package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var listEnvironmentsAppId int
var listEnvironmentsPageNumber int

var listEnvironmentsCmd = &cobra.Command{
	Use:   "list",
	Short: "List static environments",
	Long:  "List all static environments (non branch tracking, e.g. production) for the specified application.",
	Run: func(cmd *cobra.Command, args []string) {
		environmentsListUrl := fmt.Sprintf(
			"https://%s%s/environments?application_id=%s&page=%s",
			cocli.GetCoherenceDomain(),
			cocli.GetCoherenceApiPrefix(),
			fmt.Sprint(listEnvironmentsAppId),
			fmt.Sprint(listEnvironmentsPageNumber),
		)
		res, err := cocli.CoherenceApiRequest(
			"GET",
			environmentsListUrl,
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
	listEnvironmentsCmd.Flags().IntVarP(&listEnvironmentsAppId, "app_id", "a", 0, "App ID (required)")
	listEnvironmentsCmd.MarkFlagRequired("app_id")
	listEnvironmentsCmd.Flags().IntVarP(&listEnvironmentsPageNumber, "page", "p", 1, "Page number (optional - defaults to 1)")
}
