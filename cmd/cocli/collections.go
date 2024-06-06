package cocli

import (
	"github.com/spf13/cobra"
)

var collectionCmd = &cobra.Command{
	Use:   "collections",
	Short: "Coherence collection management commands",
}

func init() {
	collectionCmd.AddCommand(listCollectionsCmd)
	collectionCmd.AddCommand(getCollectionDataCmd)
}
