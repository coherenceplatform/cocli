package cocli

import (
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Coherence authentication commands",
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(printRefreshTokenCmd)
	authCmd.AddCommand(getUserInfoCmd)
}
