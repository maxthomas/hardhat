package cmd

import (
	"strconv"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/spf13/cobra"
)

var (
	// Buffered indicates if the Thrift transport will be buffered
	Buffered bool
	// Framed indicates if the Thrift transport will be framed
	Framed bool
	// BatchSize is the number of communications to send in a Store
	BatchSize int

	// StoreHost is the host of store
	StoreHost string
	// StorePort is the port of store
	StorePort uint64
)

func init() {
	RootCmd.AddCommand(LoaderCmd)
	LoaderCmd.PersistentFlags().BoolVar(&Buffered, "buffered", false, "Buffer the Thrift transport")
	LoaderCmd.PersistentFlags().BoolVar(&Framed, "framed", true, "Frame the Thrift transport")
	LoaderCmd.PersistentFlags().IntVar(&BatchSize, "batch-size", 10, "Size of communication batches in store")

	LoaderCmd.PersistentFlags().StringVar(&StoreHost, "store-host", "localhost", "Hostname of store service")
	LoaderCmd.PersistentFlags().Uint64Var(&StorePort, "store-port", 12999, "Port of store service")
}

var LoaderCmd = &cobra.Command{
	Use:   "loader",
	Short: "Put communications into a store implementation",
	Long:  `This command takes in communications and loads them into a store implementation.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	cmd.Usage()
	// },
}

// HostPort returns a string suitable for use in thrift.NewTSocket
func HostPort() string {
	return StoreHost + ":" + strconv.FormatUint(StorePort, 10)
}

// NewSocketServer uses the specified StoreHost and StorePort
// and returns a thrift.TSocket pointer, error tuple
func NewSocketServer() (*thrift.TSocket, error) {
	return thrift.NewTSocket(HostPort())
}
