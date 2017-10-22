package cmd

import "github.com/spf13/cobra"

var (
	// Buffered indicates if the Thrift transport will be buffered
	Buffered bool
	// Framed indicates if the Thrift transport will be framed
	Framed bool
)

func init() {
	RootCmd.PersistentFlags().BoolVar(&Buffered, "buffered", false, "Buffer the Thrift transport")
	RootCmd.PersistentFlags().BoolVar(&Framed, "framed", true, "Frame the Thrift transport")
}

var RootCmd = &cobra.Command{
	Use:   "hardhat",
	Short: "Hardhat is a collection of command line tools for concrete",
	Long:  `TBD`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	cmd.Usage()
	// },
}
