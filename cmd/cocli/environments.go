package cocli

import (
	"github.com/spf13/cobra"
)

var environmentCmd = &cobra.Command{
	Use:   "environments",
	Short: "Coherence environment management commands",
}

func init() {
	environmentCmd.AddCommand(deployCmd)
	environmentCmd.AddCommand(createEnvironmentCmd)
	environmentCmd.AddCommand(archiveEnvironmentCmd)
	environmentCmd.AddCommand(editEnvironmentCmd)
	environmentCmd.AddCommand(getEnvironmentCmd)
}