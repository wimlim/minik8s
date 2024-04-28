package cmd

import(
	"fmt"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete resources by filenames, stdin, resources and names",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("delete")
	},
}