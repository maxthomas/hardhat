package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// Buffered indicates if the Thrift transport will be buffered
	Buffered bool
	// Framed indicates if the Thrift transport will be framed
	Framed bool

	// Host is the server's host
	Host string

	// Port is the server's port
	Port uint32

	// ErrorRate is the percentage of errors returned
	ErrorRate float32

	// a default logger
	logger, _ = zap.NewProduction()
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVar(&Buffered, "buffered", false, "Buffer the Thrift transport")
	RootCmd.PersistentFlags().BoolVar(&Framed, "framed", true, "Frame the Thrift transport")
	RootCmd.PersistentFlags().StringVar(&Host, "host", "", "hostname to listen on")
	RootCmd.PersistentFlags().Uint32VarP(&Port, "port", "p", 41170, "port to listen on")

	RootCmd.PersistentFlags().Float32VarP(&ErrorRate, "error-rate", "e", 0.0, "percentage of errors to return")

	viper.BindPFlag("host", RootCmd.PersistentFlags().Lookup("host"))
	viper.SetDefault("host", "")
	viper.SetDefault("error_rate", 0.0)
}

func initConfig() {
	viper.SetEnvPrefix("servers")
	viper.BindEnv("host")
	viper.BindEnv("transport_buffered")
	viper.BindEnv("transport_framed")
	viper.BindEnv("error_rate")
}

var RootCmd = &cobra.Command{
	Use:   "servers",
	Short: "servers is a collection of command line tools for launching concrete servers",
	Long:  `TBD`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	cmd.Usage()
	// },
}
