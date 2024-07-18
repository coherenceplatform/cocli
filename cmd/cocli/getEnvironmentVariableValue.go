package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var getEnvVarValueCmd = &cobra.Command{
	Use:   "get-value <environment_item_id>",
	Short: "Get the value of an environment variable",
	Long:  "Retrieve the value of a specific environment variable by its ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environmentItemID := args[0]

		getValueURL := fmt.Sprintf(
			"%s/api/v1/environment_items/%s/value",
			cocli.GetCoherenceDomain(),
			environmentItemID,
		)

		res, err := cocli.CoherenceApiRequest(
			"GET",
			getValueURL,
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
