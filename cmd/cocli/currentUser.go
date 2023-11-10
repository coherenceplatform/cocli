package cocli

import (
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var currentUserCmd = &cobra.Command{
	Use:   "current_user",
	Short: "List authenticated coherence user information",
	Run: func(cmd *cobra.Command, args []string) {
		currentUserUrl := fmt.Sprintf(
			"https://%s%s/current_user",
			cocli.GetCoherenceDomain(),
			cocli.GetCoherenceApiPrefix(),
		)
		res, err := cocli.CoherenceApiRequest(
			"GET",
			currentUserUrl,
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

		fmt.Println(cocli.FormatJSONOutput(bodyBytes))
	},
}
