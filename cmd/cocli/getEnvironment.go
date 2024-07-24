package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var getEnvironmentCmd = &cobra.Command{
	Use:   "get <environment_id>",
	Short: "Get environment details",
	Long:  "Fetch and display details of a specific environment.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environmentID := args[0]

		getURL := fmt.Sprintf("%s/api/v1/environments/%s", cocli.GetCoherenceDomain(), environmentID)

		res, err := cocli.CoherenceApiRequest(
			"GET",
			getURL,
			nil,
		)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		fmt.Println(cocli.FormatJSONOutput(bodyBytes))
	},
}