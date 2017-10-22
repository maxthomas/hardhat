package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "hardhat",
	Short: "Hardhat is a collection of command line tools for concrete",
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
