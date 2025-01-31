package cocli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)


var cncCollectionId int
var printOnly bool

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

		var responseData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
			panic(err)
		}

		data, ok := responseData["data"]
		if !ok {
			panic("response does not contain 'data' key")
		}

		yamlData, err := yaml.Marshal(data)
		if err != nil {
			panic(err)
		}

		if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
			panic(err)
		}

		if printOnly {
			fmt.Println(string(yamlData))
			if err := os.Remove(outputFile); err != nil {
				fmt.Printf("Error deleting file %s: %v\n", outputFile, err)
				os.Exit(1)
			}
			os.Exit(0)
		}

		env := os.Environ()
		env = append(
			env,
			fmt.Sprintf("CNC_CONFIG_PATH=%s", outputFile),
			fmt.Sprintf("CNC_ENVIRONMENTS_PATH=%s", outputFile),
		)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT)
		done := make(chan bool, 1)

		go func() {
			cmdExec := exec.Command("cnc", args...)
			cmdExec.Stdin = os.Stdin
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			cmdExec.Env = env

			if err := cmdExec.Run(); err != nil {
				fmt.Printf("Error running command: %v\n", err)
				os.Exit(1)
			}

			done <- true
		}()

		select {
		case <-done:
			if err := os.Remove(outputFile); err != nil {
				fmt.Printf("Error deleting file %s: %v\n", outputFile, err)
				os.Exit(1)
			}
			os.Exit(0)
		case <-sigChan:
			fmt.Println("\nInterrupt received. Waiting for cleanup...")
			time.Sleep(3 * time.Second)
			fmt.Println("Exiting.")

			if err := os.Remove(outputFile); err != nil {
				fmt.Printf("Error deleting file %s: %v\n", outputFile, err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	},
}

func init() {
	cncCmd.Flags().IntVarP(&cncCollectionId, "collection_id", "c", 0, "Collection ID (required)")
	cncCmd.Flags().BoolVar(&printOnly, "print-only", false, "Print config and exit without running cnc")
	cncCmd.MarkFlagRequired("collection_id")
}