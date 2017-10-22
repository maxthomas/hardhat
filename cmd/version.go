package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hardhat",
	Long:  `All software has versions. This is Hardhat's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hardhat command line tools -- HEAD")
	},
}
