package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var upsertEnvVarsCmd = &cobra.Command{
	Use:   "upsert",
	Short: "Upsert environment variables",
	Long:  "Bulk upsert environment variables from a .env file or a string representation.",
	Run: func(cmd *cobra.Command, args []string) {
		collectionID, _ := cmd.Flags().GetInt("collection_id")
		environmentID, _ := cmd.Flags().GetInt("environment_id")
		serviceName, _ := cmd.Flags().GetString("service_name")
		envType, _ := cmd.Flags().GetString("type")
		envFile, _ := cmd.Flags().GetString("file")
		envString, _ := cmd.Flags().GetString("env-string")

		if collectionID == 0 && environmentID == 0 {
			fmt.Println("Error: either collection_id or environment_id is required")
			return
		}

		var items string
		var err error

		if envFile != "" {
			itemsBytes, err := os.ReadFile(envFile)
			if err != nil {
				fmt.Println("Error reading .env file:", err)
				return
			}
			items = string(itemsBytes)
		} else if envString != "" {
			items = envString
		} else {
			fmt.Println("Error: either --file or --env-string must be provided")
			return
		}

		payload := map[string]interface{}{
			"items": items,
			"type":  envType,
		}

		if collectionID != 0 {
			payload["collection_id"] = collectionID
		}
		if environmentID != 0 {
			payload["environment_id"] = environmentID
		}
		if serviceName != "" {
			payload["service_name"] = serviceName
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error creating JSON payload:", err)
			return
		}

		upsertURL := fmt.Sprintf("%s/api/v1/environment_items/bulk", cocli.GetCoherenceDomain())

		res, err := cocli.CoherenceApiRequest(
			"POST",
			upsertURL,
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
	upsertEnvVarsCmd.Flags().Int("collection_id", 0, "ID of the collection")
	upsertEnvVarsCmd.Flags().Int("environment_id", 0, "ID of the environment")
	upsertEnvVarsCmd.Flags().String("service_name", "", "Name of the service")
	upsertEnvVarsCmd.Flags().String("type", "standard", "Type of the environment variables (standard or secret)")
	upsertEnvVarsCmd.Flags().String("file", "", "Path to the .env file")
	upsertEnvVarsCmd.Flags().String("env-string", "", "String representation of environment variables")
}
