package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var currentUserCmd = &cobra.Command{
	Use:   "current_user",
	Short: "List authenticated user information",
	Run: func(cmd *cobra.Command, args []string) {
		// check if any creds present
		if !(cocli.CredsFileExists() || cocli.IsRefreshTokenVarSet()) {
			cocli.NotifyAuthRequired()
			return
		}

		baseDomain := fmt.Sprintf("https://%s/api/public/cli/v1", cocli.GetCoherenceDomain())
		res, err := cocli.CoherenceApiRequest(
			"GET",
			fmt.Sprintf("%s/current_user", baseDomain),
			nil,
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		if res.StatusCode == 401 {
			// handle unauthorized
		}

		fmt.Println(string(bodyBytes))
	},
}
