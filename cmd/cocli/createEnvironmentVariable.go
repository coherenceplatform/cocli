package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var createEnvVarCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment variable",
	Long:  "Create a new environment variable for a specific environment or collection.",
	Run: func(cmd *cobra.Command, args []string) {
		targetEnvName, _ := cmd.Flags().GetString("target_env_name")
		itemType, _ := cmd.Flags().GetString("type")
		collectionID, _ := cmd.Flags().GetInt("collection_id")
		environmentID, _ := cmd.Flags().GetInt("environment_id")
		serviceName, _ := cmd.Flags().GetString("service_name")
		value, _ := cmd.Flags().GetString("value")
		secretID, _ := cmd.Flags().GetString("secret_id")
		outputName, _ := cmd.Flags().GetString("output_name")
		alias, _ := cmd.Flags().GetString("alias")

		if targetEnvName == "" {
			fmt.Println("Error: target_env_name is required")
			return
		}

		if collectionID == 0 && environmentID == 0 {
			fmt.Println("Error: either collection_id or environment_id is required")
			return
		}

		payload := map[string]interface{}{
			"target_env_name": targetEnvName,
			"type":            itemType,
			"service_name":    serviceName,
			"value":           value,
			"secret_id":       secretID,
			"output_name":     outputName,
			"alias":           alias,
		}

		if collectionID != 0 {
			payload["collection_id"] = collectionID
		}
		if environmentID != 0 {
			payload["environment_id"] = environmentID
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error creating JSON payload:", err)
			return
		}

		createURL := fmt.Sprintf("%s/api/v1/environment_items", cocli.GetCoherenceDomain())

		res, err := cocli.CoherenceApiRequest(
			"POST",
			createURL,
			bytes.NewBuffer(jsonPayload),
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

func init() {
	createEnvVarCmd.Flags().String("target_env_name", "", "Name of the environment variable (required)")
	createEnvVarCmd.Flags().String("type", "standard", "Type of the environment variable (standard, secret, output, alias)")
	createEnvVarCmd.Flags().Int("collection_id", 0, "ID of the collection")
	createEnvVarCmd.Flags().Int("environment_id", 0, "ID of the environment")
	createEnvVarCmd.Flags().String("service_name", "", "Name of the service")
	createEnvVarCmd.Flags().String("value", "", "Value of the environment variable")
	createEnvVarCmd.Flags().String("secret_id", "", "ID of the secret (for secret type)")
	createEnvVarCmd.Flags().String("output_name", "", "Name of the output (for output type)")
	createEnvVarCmd.Flags().String("alias", "", "Alias (for alias type)")

	createEnvVarCmd.MarkFlagRequired("target_env_name")
}