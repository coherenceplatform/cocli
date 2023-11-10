package cocli

import (
	"fmt"

	cocli "github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var printRefreshTokenCmd = &cobra.Command{
	Use:   "print_refresh_token",
	Short: "Print your Coherence API refresh token to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		token := cocli.GetTokenFromCredsFile()

		fmt.Printf("COCLI_REFRESH_TOKEN='%s'\n", token.RefreshToken)
	},
}
