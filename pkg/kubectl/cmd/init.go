package cmd

import(
	"fmt"
	"github.com/spf13/cobra"
)

var AllCommands = cobra.Command{
	Use: "kubectl",
	Short: "Kubernetes command line tool",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println(cmd.UsageString())
	},
}

func init(){
	AllCommands.AddCommand(applyCmd)
	AllCommands.AddCommand(deleteCmd)
}

func Execute(){
	AllCommands.Execute()
}