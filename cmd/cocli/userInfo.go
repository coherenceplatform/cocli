package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var getUserInfoCmd = &cobra.Command{
	Use:   "userinfo",
	Short: "List authenticated user idtoken information",
	Run: func(cmd *cobra.Command, args []string) {
		// check if any creds present
		if !(cocli.CredsFileExists() || cocli.IsRefreshTokenVarSet()) {
			cocli.NotifyAuthRequired()
			return
		}

		baseDomain := fmt.Sprintf("https://%s", cocli.GetAuthDomain())
		res, err := cocli.OauthProviderApiRequest(
			"GET",
			fmt.Sprintf("%s/userinfo", baseDomain),
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
		// respObject := make(map[string]string)
		// json.Unmarshal(bodyBytes, &respObject)
		// fmt.Println(respObject)
	},
}
