package cmd

import (
	"fmt"
	"os"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
	"github.com/spf13/cobra"
)

const (
	summaryTypesURL string = "http://hltcoe.github.io/concrete/schema/summarization.html#Enum_SummarySourceType"
)

var (
	// SummarizeHost is the host of store
	SummarizeHost string
	// SummarizePort is the port of store
	SummarizePort uint64
	// UUIDs are uuids to add to the summarization request
	UUIDs []string
	// TypeString is the type of summarization request
	TypeString string

	// MaxTokens is the maximum number of tokens to submit
	MaxTokens uint
)

func init() {
	RootCmd.AddCommand(SummarizeCmd)

	SummarizeCmd.PersistentFlags().StringVar(&SummarizeHost, "summ-host", "localhost", "Hostname of summarization service")
	SummarizeCmd.PersistentFlags().Uint64Var(&SummarizePort, "summ-port", 22999, "Port of summarization service")
	SummarizeCmd.Flags().StringVarP(&TypeString, "type", "t", "entity", "Summary source type (example: entity)")
	SummarizeCmd.Flags().StringSliceVar(&UUIDs, "uuids", []string{}, "UUIDs to query")
	SummarizeCmd.Flags().UintVar(&MaxTokens, "max-tokens", 250, "maximum tokens in summary")
}

var SummarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Client for interacting with a summarize service",
	Long:  `This provides a client for interacting with a summarization service.`,
	Run: func(cmd *cobra.Command, args []string) {
		summType, err := goncrete.SummarySourceTypeFromString(strings.ToUpper(TypeString))
		if err != nil {
			fmt.Printf("error parsing summary source type: %v\n", err)
			fmt.Println("check available types: " + summaryTypesURL)
			os.Exit(-1)
		}

		var transportFactory thrift.TTransportFactory
		if Buffered {
			transportFactory = thrift.NewTBufferedTransportFactory(8192)
		} else {
			transportFactory = thrift.NewTTransportFactory()
		}

		if Framed {
			transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
		}
		protocolFactory := thrift.NewTCompactProtocolFactory()
		tport, err := NewSocketServer(SummarizeHost, SummarizePort)
		if err != nil {
			fmt.Printf("error creating socket: %v\n", err)
			os.Exit(-1)
		}

		wrapped := transportFactory.GetTransport(tport)
		defer wrapped.Close()
		if err = wrapped.Open(); err != nil {
			fmt.Printf("error connecting to store: %v\n", err)
			os.Exit(-1)
		}

		summCli := goncrete.NewSummarizationServiceClientFactory(wrapped, protocolFactory)
		alive, err := summCli.Alive()
		if err != nil || !alive {
			fmt.Printf("server not alive / error querying alive: %v", err)
			os.Exit(-1)
		}

		summReq := goncrete.NewSummarizationRequest()
		thriftMaxTokens := int32(MaxTokens)
		summReq.MaximumTokens = thrift.Int32Ptr(thriftMaxTokens)
		for _, uuid := range UUIDs {
			concUUID := goncrete.NewUUID()
			concUUID.UuidString = uuid
			summReq.SourceIds = append(summReq.SourceIds, concUUID)
		}
		summReq.SourceType = goncrete.SummarySourceTypePtr(summType)
		resp, err := summCli.Summarize(summReq)
		if err != nil {
			fmt.Printf("error summary call: %v\n", err)
			os.Exit(-1)
		}
		if resp == nil {
			fmt.Println("received nil Summary")
		} else if resp.IsSetSummaryCommunication() && resp.GetSummaryCommunication().IsSetText() {
			fmt.Println(resp.GetSummaryCommunication().GetText())
		} else {
			fmt.Println("no summary communication text received")
		}
	},
}
