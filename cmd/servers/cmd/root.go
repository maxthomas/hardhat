package cmd

import (
	"fmt"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"
	thriftutil "github.com/maxthomas/hardhat/pkg/util/thrift"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// viper consts
	prefix = "hardhat"

	// viper env vars
	host     = "host"
	port     = "port"
	buffered = "transport_buffered"
	framed   = "transport_framed"
)

var (
	// a default logger
	logger, _ = zap.NewProduction()

	// ensure that config impls ObjectMarshaler
	_ zapcore.ObjectMarshaler = (*config)(nil)
)

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault(host, "")
	viper.SetDefault(port, 41170)
	viper.SetDefault(buffered, false)
	viper.SetDefault(framed, true)
}

func initConfig() {
	viper.SetEnvPrefix(prefix)

	viper.BindEnv(host)
	viper.BindEnv(port)
	viper.BindEnv(buffered)
	viper.BindEnv(framed)
}

type config struct {
	host     string
	port     int
	buffered bool
	framed   bool
}

func getConfig() config {
	p := viper.GetInt(port)
	if p == 0 {
		fmt.Printf("didn't parse port correctly; is this a port? %v\n", viper.GetString(port))
		os.Exit(1)
	}

	return config{
		viper.GetString(host),
		p,
		viper.GetBool(buffered),
		viper.GetBool(framed),
	}
}

func (c *config) transportFactory() thrift.TTransportFactory {
	return thriftutil.TransportFactory(c.buffered, c.framed)
}

func (c *config) serverSocket() (*thrift.TServerSocket, error) {
	return thriftutil.NewServerSocket(c.host, uint64(c.port))
}

func (c *config) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddString("host", c.host)
	oe.AddInt("port", c.port)
	oe.AddBool("buffered", c.buffered)
	oe.AddBool("framed", c.framed)
	return nil
}

var RootCmd = &cobra.Command{
	Use:   "servers",
	Short: "servers is a collection of command line tools for launching concrete servers",
	Long:  `TBD`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	cmd.Usage()
	// },
}
