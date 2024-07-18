package cocli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var listEnvVarsCmd = &cobra.Command{
	Use:   "list",
	Short: "List environment variables",
	Long:  "List all environment variables for a specific environment or collection.",
	Run: func(cmd *cobra.Command, args []string) {
		environmentID, _ := cmd.Flags().GetInt("environment_id")
		collectionID, _ := cmd.Flags().GetInt("collection_id")

		if environmentID == 0 && collectionID == 0 {
			fmt.Println("Error: Either environment_id or collection_id must be provided")
			return
		}

		if environmentID != 0 && collectionID != 0 {
			fmt.Println("Error: Please provide either environment_id or collection_id, not both")
			return
		}

		var url string
		if environmentID != 0 {
			url = fmt.Sprintf("%s/api/v1/environments/%d", cocli.GetCoherenceDomain(), environmentID)
		} else {
			url = fmt.Sprintf("%s/api/v1/collections/%d", cocli.GetCoherenceDomain(), collectionID)
		}

		res, err := cocli.CoherenceApiRequest("GET", url, nil)
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

		var responseData map[string]interface{}
		err = json.Unmarshal(bodyBytes, &responseData)
		if err != nil {
			fmt.Println("Error parsing JSON response:", err)
			return
		}

		output := make(map[string]interface{})

		if envItems, ok := responseData["environment_items"]; ok {
			output["environment_items"] = envItems
		}

		if managedEnvItems, ok := responseData["managed_environment_items"]; ok {
			output["managed_environment_items"] = managedEnvItems
		}

		outputJSON, err := json.Marshal(output)
		if err != nil {
			fmt.Println("Error creating JSON output:", err)
			return
		}

		fmt.Println(cocli.FormatJSONOutput(outputJSON))
	},
}

func init() {
	listEnvVarsCmd.Flags().Int("environment_id", 0, "ID of the environment")
	listEnvVarsCmd.Flags().Int("collection_id", 0, "ID of the collection")
}
