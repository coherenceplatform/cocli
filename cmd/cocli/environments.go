package cocli

import (
	"github.com/spf13/cobra"
)

var environmentCmd = &cobra.Command{
	Use:   "environments",
	Short: "Coherence environment management commands",
}

func init() {
	environmentCmd.AddCommand(deployEnvironmentCmd)
	environmentCmd.AddCommand(createEnvironmentCmd)
	environmentCmd.AddCommand(editEnvironmentCmd)
	environmentCmd.AddCommand(getEnvironmentCmd)
}