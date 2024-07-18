package cocli

import (
	"github.com/spf13/cobra"
)

var envVarCmd = &cobra.Command{
	Use:   "env-vars",
	Short: "Coherence environment variable management commands",
}

func init() {
	envVarCmd.AddCommand(listEnvVarsCmd)
	envVarCmd.AddCommand(createEnvVarCmd)
	envVarCmd.AddCommand(deleteEnvVarCmd)
	envVarCmd.AddCommand(upsertEnvVarsCmd)
	envVarCmd.AddCommand(getEnvVarValueCmd)
}
