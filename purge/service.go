package purge

import (
	"github.com/spf13/cobra"
)

var cmdSvc = &cobra.Command{
	Use: "service",
	Short: "Purge stopped services",
	Long: "Purge stopped services",
	RunE: RunCmdService,
}

func RunCmdService(cmd *cobra.Command, args []string) error {
	return RemoveServices()
}


func RemoveServices(filters ...Filter) error {
	return nil
}
