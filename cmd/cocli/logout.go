package cocli

import (
	"fmt"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored Coherence API credentials",
	Long: `
	Clears stored Coherence API credentials.
	This does not clear credentials stored via environment variables.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cocli.ClearCredentialsFile() {
			fmt.Println("Successfully logged out.")
		}
	},
}
