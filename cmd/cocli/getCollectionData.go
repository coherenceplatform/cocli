package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var getCollectionDataId int

var getCollectionDataCmd = &cobra.Command{
	Use:   "cnc_data",
	Short: "Get cnc formatted collection data",
	Run: func(cmd *cobra.Command, args []string) {
		collectionDataUrl := fmt.Sprintf(
			"%s/api/v1/collections/%s/cnc_config",
			cocli.GetCoherenceDomain(),
			fmt.Sprint(getCollectionDataId),
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

		fmt.Println(cocli.FormatJSONOutput(bodyBytes))
	},
}

func init() {
	getCollectionDataCmd.Flags().IntVarP(&getCollectionDataId, "collection_id", "c", 0, "Collection ID (required)")
	getCollectionDataCmd.MarkFlagRequired("collection_id")
}
