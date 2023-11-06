package cocli

import (
	"github.com/spf13/cobra"
)

var applicationCmd = &cobra.Command{
	Use:   "apps",
	Short: "Coherence application management commands",
}

func init() {
	applicationCmd.AddCommand(listAppsCmd)
}
