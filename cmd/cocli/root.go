package cocli

import (
	"fmt"
	"os"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "cocli",
	Short:   "cocli - A cli for interacting with the Coherence API",
	Version: cocli.GetCliVersion(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops, there was an error while executing your command '%s'", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(featureCmd)
	rootCmd.AddCommand(currentUserCmd)
	rootCmd.AddCommand(applicationCmd)
	rootCmd.AddCommand(environmentCmd)
	cocli.RunCliVersionCheck()
}
