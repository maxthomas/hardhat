package cmd

import (
	"fmt"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
	"github.com/maxthomas/hardhat/pkg/summarizer"
	thriftutil "github.com/maxthomas/hardhat/pkg/util/thrift"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ()

var SummarizerCmd = &cobra.Command{
	Use:   "summarizer",
	Short: "Stand up an summarize server",
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := thriftutil.NewServerSocket(Host, uint64(Port))
		if err != nil {
			fmt.Printf("error open socket: %v\n", err.Error())
			os.Exit(3)
		}
		defer srv.Close()

		summ := summarizer.NewCopier(logger)
		transFactory := thriftutil.TransportFactory(Buffered, Framed)
		protoFactory := goncrete.DefaultProtocolFactory()
		proc := goncrete.NewSummarizationServiceProcessor(summ)
		srvr := thrift.NewTSimpleServer4(proc, srv, transFactory, protoFactory)

		logger.Info("server starting")
		if err := srvr.Serve(); err != nil {
			logger.Info("server exiting", zap.Error(err))
		}
	},
}

func init() {
	RootCmd.AddCommand(SummarizerCmd)
}
