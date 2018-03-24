package app

import (
	"strconv"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/spf13/cobra"
)

var (
	// StoreHost is the host of store
	StoreHost string
	// StorePort is the port of store
	StorePort uint64
)

func init() {
	RootCmd.AddCommand(StoreCmd)

	StoreCmd.PersistentFlags().StringVar(&StoreHost, "store-host", "localhost", "Hostname of store service")
	StoreCmd.PersistentFlags().Uint64Var(&StorePort, "store-port", 12999, "Port of store service")
}

var StoreCmd = &cobra.Command{
	Use:   "store",
	Short: "Put communications into a store implementation",
	Long:  `This command takes in communications and loads them into a store implementation.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	cmd.Usage()
	// },
}

// HostPort returns a string suitable for use in thrift.NewTSocket
func HostPort(host string, port uint64) string {
	return host + ":" + strconv.FormatUint(port, 10)
}

// NewSocketServer uses the specified StoreHost and StorePort
// and returns a thrift.TSocket pointer, error tuple
func NewSocketServer(host string, port uint64) (*thrift.TSocket, error) {
	return thrift.NewTSocket(HostPort(host, port))
}
