package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var archiveFeatureId int

var archiveFeatureCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a coherence feature",
	Run: func(cmd *cobra.Command, args []string) {
		archiveFeaturePayload := map[string]string{
			"status": "archived",
		}
		payloadBytes, err := json.Marshal(archiveFeaturePayload)
		if err != nil {
			panic(err)
		}
		payload := bytes.NewBuffer(payloadBytes)

		featuresUpdateUrl := fmt.Sprintf(
			"https://%s%s/features/%s",
			cocli.GetCoherenceDomain(),
			cocli.GetCoherenceApiPrefix(),
			fmt.Sprint(archiveFeatureId),
		)
		res, err := cocli.CoherenceApiRequest(
			"PUT",
			featuresUpdateUrl,
			payload,
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(cocli.FormatJSONOutput(bodyBytes))
	},
}

func init() {
	archiveFeatureCmd.Flags().IntVarP(&archiveFeatureId, "feature_id", "f", 0, "Feature ID (required)")
	archiveFeatureCmd.MarkFlagRequired("feature_id")
}
