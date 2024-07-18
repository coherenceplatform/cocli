package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var deleteEnvVarCmd = &cobra.Command{
	Use:   "delete <environment_item_id>",
	Short: "Delete an environment variable",
	Long:  "Delete an environment variable by its ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environmentItemID := args[0]

		deleteURL := fmt.Sprintf(
			"%s/api/v1/environment_items/%s",
			cocli.GetCoherenceDomain(),
			environmentItemID,
		)

		res, err := cocli.CoherenceApiRequest(
			"DELETE",
			deleteURL,
			nil,
		)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode == 204 {
			fmt.Printf("Environment variable with ID %s has been successfully deleted.\n", environmentItemID)
			return
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		fmt.Println(cocli.FormatJSONOutput(bodyBytes))
	},
}
