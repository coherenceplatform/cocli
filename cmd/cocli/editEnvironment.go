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

var editEnvironmentCmd = &cobra.Command{
	Use:   "edit <environment_id>",
	Short: "Edit an existing environment",
	Long:  "Edit properties of an existing environment.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environmentID := args[0]

		name, _ := cmd.Flags().GetString("name")
		deploymentType, _ := cmd.Flags().GetString("type")
		status, _ := cmd.Flags().GetString("status")
		starred, _ := cmd.Flags().GetBool("starred")
		customDomains, _ := cmd.Flags().GetStringSlice("custom_domains")
		customSSLCertificates, _ := cmd.Flags().GetStringSlice("custom_ssl_certificates")
		cloneServicesFromEnvironmentID, _ := cmd.Flags().GetInt("clone_services_from_environment_id")
		cloneServiceInfos, _ := cmd.Flags().GetStringToString("clone_service_infos")

		payload := make(map[string]interface{})

		if name != "" {
			payload["name"] = name
		}
		if deploymentType != "" {
			payload["type"] = deploymentType
		}
		if status != "" {
			payload["status"] = status
		}
		if cmd.Flags().Changed("starred") {
			payload["starred"] = starred
		}
		if len(customDomains) > 0 {
			payload["custom_domains"] = customDomains
		}
		if len(customSSLCertificates) > 0 {
			payload["custom_ssl_certificates"] = customSSLCertificates
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

		editURL := fmt.Sprintf("%s/api/v1/environments/%s", cocli.GetCoherenceDomain(), environmentID)

		res, err := cocli.CoherenceApiRequest(
			"PUT",
			editURL,
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
	editEnvironmentCmd.Flags().String("name", "", "New name for the environment")
	editEnvironmentCmd.Flags().String("type", "", "New deployment type for the environment")
	editEnvironmentCmd.Flags().String("status", "", "New status for the environment")
	editEnvironmentCmd.Flags().Bool("starred", false, "Whether the environment should be starred")
	editEnvironmentCmd.Flags().StringSlice("custom_domains", []string{}, "Custom domains for the environment")
	editEnvironmentCmd.Flags().StringSlice("custom_ssl_certificates", []string{}, "Custom SSL certificates for the environment")
	editEnvironmentCmd.Flags().Int("clone_services_from_environment_id", 0, "ID of the environment to clone services from")
	editEnvironmentCmd.Flags().StringToString("clone_service_infos", nil, "Service info for cloning (format: name=branch_name)")
}