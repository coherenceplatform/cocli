package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var listAppsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all coherence applications",
	Long:  "List all coherence applications that are accessible by the currently authenticated user.",
	Run: func(cmd *cobra.Command, args []string) {
		appsListUrl := fmt.Sprintf(
			"https://%s%s/applications",
			cocli.GetCoherenceDomain(),
			cocli.GetCoherenceApiPrefix(),
		)
		res, err := cocli.CoherenceApiRequest(
			"GET",
			appsListUrl,
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
