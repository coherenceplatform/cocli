package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

type EnvironmentServiceDeploy struct {
	Name       string `json:"name"`
	BranchName string `json:"branch_name,omitempty"`
	CommitSHA  string `json:"commit_sha,omitempty"`
}

type EnvironmentDeploy struct {
	Services       []EnvironmentServiceDeploy `json:"services"`
	ConfigureInfra bool                       `json:"configure_infra,omitempty"`
}

var deployEnvironmentCmd = &cobra.Command{
	Use:   "deploy <environment_id> <json_file_or_string>",
	Short: "Deploy to a specific environment",
	Long:  "Deploy services to a specific Coherence environment using a JSON configuration.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		environmentID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid environment ID")
			return
		}

		var deployData EnvironmentDeploy
		jsonInput := args[1]

		// Check if the input is a file
		if _, err := os.Stat(jsonInput); err == nil {
			jsonFile, err := os.ReadFile(jsonInput)
			if err != nil {
				fmt.Println("Error reading JSON file:", err)
				return
			}
			err = json.Unmarshal(jsonFile, &deployData)
		} else {
			// Treat input as a JSON string
			err = json.Unmarshal([]byte(jsonInput), &deployData)
		}

		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		if len(deployData.Services) == 0 {
			fmt.Println("At least one service must be specified")
			return
		}

		jsonData, err := json.Marshal(deployData)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		deployURL := fmt.Sprintf(
			"%s/api/v1/environments/%d/deploy",
			cocli.GetCoherenceDomain(),
			environmentID,
		)

		res, err := cocli.CoherenceApiRequest(
			"PUT",
			deployURL,
			bytes.NewBuffer(jsonData),
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
