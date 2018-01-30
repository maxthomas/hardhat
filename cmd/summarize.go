package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
	"github.com/spf13/cobra"

	"github.com/maxthomas/hardhat/utils"
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

	// IsCapabilities calls GetCapabilities then quits
	IsCapabilities bool

	// CommunicationPath takes in a path to a communication to
	// send to the summary service
	CommunicationPath string
)

func init() {
	RootCmd.AddCommand(SummarizeCmd)

	SummarizeCmd.PersistentFlags().StringVar(&SummarizeHost, "summ-host", "localhost", "Hostname of summarization service")
	SummarizeCmd.PersistentFlags().Uint64Var(&SummarizePort, "summ-port", 22999, "Port of summarization service")
	SummarizeCmd.Flags().StringVarP(&TypeString, "type", "t", "entity", "Summary source type (example: entity)")
	SummarizeCmd.Flags().StringSliceVar(&UUIDs, "uuids", []string{}, "UUIDs to query")
	SummarizeCmd.Flags().UintVar(&MaxTokens, "max-tokens", 250, "maximum tokens in summary")
	SummarizeCmd.Flags().BoolVar(&IsCapabilities, "capabilities", false, "perform a getCapabilities call then exit")
	SummarizeCmd.Flags().StringVar(&CommunicationPath, "path", "", "path to concrete communication to send")

}

var SummarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Client for interacting with a summarize service",
	Long:  `This provides a client for interacting with a summarization service.`,
	Run: func(cmd *cobra.Command, args []string) {
		var sourceComm *goncrete.Communication
		var summType goncrete.SummarySourceType
		var err error

		// if this is not a capabilities check, validate the input
		if !IsCapabilities {
			summType, err = goncrete.SummarySourceTypeFromString(strings.ToUpper(TypeString))
			if err != nil {
				fmt.Printf("error parsing summary source type: %v\n", err)
				fmt.Println("check available types: " + summaryTypesURL)
				os.Exit(-1)
			}
		}

		// if the CommunicationPath flag is set, validate that it can
		// be successfully deserialized
		if CommunicationPath != "" {
			if _, err := os.Stat(CommunicationPath); err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("no file at: %v\n", CommunicationPath)
					os.Exit(129)
				} else {
					fmt.Printf("error occurred during file check: %v\n", err.Error())
					os.Exit(130)
				}
			}
			commFile, err := os.Open(CommunicationPath)
			if err != nil {
				fmt.Printf("error loading file: %v\n", err.Error())
				os.Exit(131)
			}
			defer commFile.Close()
			commBytes, err := ioutil.ReadAll(commFile)
			if err != nil {
				fmt.Printf("error reading file contents: %v\n", err.Error())
				os.Exit(132)
			}
			sourceComm, err = utils.LoadCommunication(commBytes)
			if err != nil {
				fmt.Printf("error deserializing communication: %v\n", err.Error())
				os.Exit(133)
			}
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
		if IsCapabilities {
			caps, capErr := summCli.GetCapabilities()
			if capErr != nil {
				fmt.Printf("error during get capabilities call: %v", capErr.Error())
				os.Exit(-1)
			}

			fmt.Printf("Received %v capabilities from server:\n", len(caps))
			for idx := range caps {
				cap := caps[idx]
				fmt.Printf("\t%v\n", cap.String())
			}
			return
		}

		summReq := goncrete.NewSummarizationRequest()
		if sourceComm != nil {
			fmt.Printf("setting source communication: %v\n", sourceComm.GetID())
			summReq.SourceCommunication = sourceComm
		}

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
			fmt.Printf("received summary text: %v\n", resp.GetSummaryCommunication().GetText())
		} else {
			fmt.Println("no summary communication text received")
		}
	},
}
