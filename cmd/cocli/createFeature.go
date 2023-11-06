package cocli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createFeatureAppId int

var createFeatureCmd = &cobra.Command{
	Use:   "create <branch_name>",
	Short: "Create a new coherence feature",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		fmt.Printf("branch %s | appId %s", branchName, fmt.Sprint(createFeatureAppId))
	},
}

func init() {
	createFeatureCmd.Flags().IntVarP(&createFeatureAppId, "app_id", "a", 0, "App ID (required)")
	createFeatureCmd.MarkFlagRequired("app_id")
}
