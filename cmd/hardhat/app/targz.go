package app

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/hltcoe/goncrete"
	utils "github.com/maxthomas/hardhat/pkg/util"
	"github.com/spf13/cobra"
)

var (
	GZIPCompressed bool
)

func init() {
	StoreCmd.AddCommand(loadTarGzCmd)
	loadTarGzCmd.Flags().BoolVarP(&GZIPCompressed, "gzipped", "z", true, "Is the tar archive gzip compressed?")
}

var loadTarGzCmd = &cobra.Command{
	Use:   "targz",
	Short: "Reads a .tar or .tar.gz file of communications into a store implementation",
	Long: `This command reads in a gzip compressed tar file (or normal tar file)
of communications and inserts them into a store implementation.

It requires a single position argument, the path to a tar archive.

By default, the tar file is assumed to be gzip compressed. This can be disabled.`,
	// Args: cobra.ExactArgs(1),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Exactly one argument is required: the file path")
		}

		stat, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			return fmt.Errorf("File does not exist at: %v", args[0])
		}

		if stat.IsDir() {
			return fmt.Errorf("Path is a directory: %v", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		commFile, err := os.Open(path)
		if err != nil {
			fmt.Printf("error opening file: %v\n", err.Error())
			os.Exit(-1)
		}
		defer commFile.Close()
		var rdr io.Reader = commFile
		if GZIPCompressed {
			rdr, err = gzip.NewReader(commFile)
			if err != nil {
				fmt.Printf("error opening gzip file, is it really a gzip file?: %v\n", err.Error())
				os.Exit(-2)
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
		tport, err := NewSocketServer(StoreHost, StorePort)
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

		storeImpl := goncrete.NewStoreCommunicationServiceClientFactory(wrapped, protocolFactory)
		alive, err := storeImpl.Alive()
		if err != nil || !alive {
			fmt.Printf("server not alive / error querying alive: %v", err)
			os.Exit(-1)
		}

		tarRdr := tar.NewReader(rdr)
		errorChan := make(chan error, 2)
		tarItems := make(chan []byte, 10)

		// load the file in parallel
		go func(trdr *tar.Reader, data chan<- []byte, errorC chan<- error) {
			errorChan <- utils.LoadTarFile(trdr, data)
			close(data)
			close(errorC)
			fmt.Println("Finished loading file")
		}(tarRdr, tarItems, errorChan)

		memBuf := thrift.NewTMemoryBuffer()
		deser := thrift.NewTDeserializer()
		deser.Protocol = protocolFactory.GetProtocol(memBuf)
		deser.Transport = memBuf

		for data := range tarItems {
			comm := goncrete.NewCommunication()
			if err = deser.Read(comm, data); err != nil {
				fmt.Printf("error deserializing: %v", err)
				os.Exit(-1)
			}
			if err = storeImpl.Store(comm); err != nil {
				fmt.Printf("error during store call: %v", err)
				os.Exit(-1)
			}

			memBuf.Reset()
		}
		if err = <-errorChan; err != nil {
			fmt.Printf("error during processing: %v\n", err)
			os.Exit(-1)
		}
	},
}
