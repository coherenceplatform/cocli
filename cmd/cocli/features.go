package cocli

import (
	"github.com/spf13/cobra"
)

var featureCmd = &cobra.Command{
	Use:   "features",
	Short: "Coherence feature management commands",
}

func init() {
	featureCmd.AddCommand(listFeaturesCmd)
	featureCmd.AddCommand(createFeatureCmd)
}
