package kubectl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var AllCommands = cobra.Command{
	Use:   "kubectl",
	Short: "Kubernetes command line tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

func init() {
	AllCommands.AddCommand(applyCmd)
	AllCommands.AddCommand(deleteCmd)
	AllCommands.AddCommand(getCmd)
	AllCommands.AddCommand(describeCmd)
	AllCommands.AddCommand(invokeCmd)
}

func Execute() {
	AllCommands.Execute()
}
