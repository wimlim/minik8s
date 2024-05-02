package cmd

import(
	"fmt"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use: "get",
	Short: "Display one or many resources",
	Run: getHandler,
}

func getHandler(cmd *cobra.Command, args []string){
	if(len(args) == 0){
		fmt.Println("no args")
		return
	}
	fmt.Println("get handler")
}