package cocli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var cncCollectionId int

var cncCmd = &cobra.Command{
	Use:   "cnc",
	Short: "Coherence cnc wrapper",
	Long:  "Coherence cnc wrapper - runs cnc commands using configuration from coherence",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && args[0] == "provision" {
			fmt.Printf("cnc 'provision' commands are not supported via cocli")
			os.Exit(1)
		}

		outputFile := fmt.Sprintf(
			"/tmp/.cocli_cnc_config_%d.yml",
			cncCollectionId,
		)
		collectionDataUrl := fmt.Sprintf(
			"%s/api/v1/collections/%s/cnc_config",
			cocli.GetCoherenceDomain(),
			fmt.Sprint(cncCollectionId),
		)
		res, err := cocli.CoherenceApiRequest(
			"GET",
			collectionDataUrl,
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

		// Parse the JSON response
		var responseData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
			panic(err)
		}

		// Extract the "data" key
		data, ok := responseData["data"]
		if !ok {
			panic("response does not contain 'data' key")
		}

		// Convert the extracted data to YAML
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			panic(err)
		}

		// Write the YAML to a file
		if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
			panic(err)
		}

		// Build the command with the new argument structure
		commandArgs := []string{"-e", outputFile, "-f", outputFile}
		commandArgs = append(commandArgs, args...)

		// Run a shell command using the arguments passed to the Cobra command
		cmdExec := exec.Command("cnc", commandArgs...)
		cmdExec.Stdin = os.Stdin
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			fmt.Printf("Error running command: %v\n", err)
			os.Exit(1)
		}

		if err := os.Remove(outputFile); err != nil {
			fmt.Printf("Error deleting file %s: %v\n", outputFile, err)
			os.Exit(1)
		}
	},
}

func init() {
	cncCmd.Flags().IntVarP(&cncCollectionId, "collection_id", "c", 0, "Collection ID (required)")
	cncCmd.MarkFlagRequired("collection_id")
}
