package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

type EnvironmentCloneServiceInfo struct {
	Name            string `json:"name"`
	TrackBranchName string `json:"track_branch_name,omitempty"`
}

var createEnvironmentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Long:  "Create a new environment in a specific collection.",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		collectionID, _ := cmd.Flags().GetInt("collection_id")
		deploymentType, _ := cmd.Flags().GetString("type")
		cloneEnvironmentID, _ := cmd.Flags().GetInt("clone_environment_id")
		cloneServicesFromEnvironmentID, _ := cmd.Flags().GetInt("clone_services_from_environment_id")
		cloneServiceInfos, _ := cmd.Flags().GetStringToString("clone_service_infos")

		if name == "" || collectionID == 0 {
			fmt.Println("Error: name and collection_id are required")
			return
		}

		payload := map[string]interface{}{
			"name":          name,
			"collection_id": collectionID,
		}

		if deploymentType != "" {
			payload["type"] = deploymentType
		}

		if cloneEnvironmentID != 0 {
			payload["clone_environment_id"] = cloneEnvironmentID
		}

		if cloneServicesFromEnvironmentID != 0 {
			payload["clone_services_from_environment_id"] = cloneServicesFromEnvironmentID
		}

		if len(cloneServiceInfos) > 0 {
			var serviceInfos []EnvironmentCloneServiceInfo
			for name, branchName := range cloneServiceInfos {
				serviceInfos = append(serviceInfos, EnvironmentCloneServiceInfo{
					Name:            name,
					TrackBranchName: branchName,
				})
			}
			payload["clone_service_infos"] = serviceInfos
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error creating JSON payload:", err)
			return
		}

		createURL := fmt.Sprintf("%s/api/v1/environments", cocli.GetCoherenceDomain())

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
	createEnvironmentCmd.Flags().String("name", "", "Name of the environment (required)")
	createEnvironmentCmd.Flags().Int("collection_id", 0, "ID of the collection (required)")
	createEnvironmentCmd.Flags().String("type", "ephemeral", "Deployment type of the environment")
	createEnvironmentCmd.Flags().Int("clone_environment_id", 0, "ID of the environment to clone from")
	createEnvironmentCmd.Flags().Int("clone_services_from_environment_id", 0, "ID of the environment to clone services from")
	createEnvironmentCmd.Flags().StringToString("clone_service_infos", nil, "Service info for cloning (format: name=branch_name)")

	createEnvironmentCmd.MarkFlagRequired("name")
	createEnvironmentCmd.MarkFlagRequired("collection_id")
}