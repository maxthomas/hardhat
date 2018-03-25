package cmd

import (
	"fmt"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
	"github.com/maxthomas/hardhat/pkg/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var StorerCmd = &cobra.Command{
	Use:   "storer",
	Short: "Stand up an store server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		logger.Info("initialized", zap.Object("configuration", &cfg))
		srv, err := cfg.serverSocket()
		if err != nil {
			fmt.Printf("error open socket: %v\n", err.Error())
			os.Exit(3)
		}
		defer srv.Close()

		annotator := &store.PrintStorer{}
		transFactory := cfg.transportFactory()
		protoFactory := goncrete.DefaultProtocolFactory()
		proc := goncrete.NewStoreCommunicationServiceProcessor(annotator)
		srvr := thrift.NewTSimpleServer4(proc, srv, transFactory, protoFactory)

		logger.Info("store server starting")
		if err := srvr.Serve(); err != nil {
			logger.Info("server exiting", zap.Error(err))
		}
	},
}

func init() {
	RootCmd.AddCommand(StorerCmd)
}
