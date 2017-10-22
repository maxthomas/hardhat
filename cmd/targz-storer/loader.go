package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/hltcoe/goncrete"
	flag "github.com/spf13/pflag"

	"git.apache.org/thrift.git/lib/go/thrift"

	"go.uber.org/zap"
)

var (
	log, _ = zap.NewProduction()
)

func main() {
	var (
		inputFile   = flag.String("input-file", "", "input .tar.gz file of communications")
		fetchServer = flag.String("store-server-address", "localhost:12999", "store server host:port")
	)
	flag.Parse()

	if *inputFile == "" {
		flag.Usage()
		return
	}

	transFactory := goncrete.DefaultTransportFactory()
	protoFactory := goncrete.DefaultProtocolFactory()

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal("error open input file", zap.Error(err))
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal("failed to create gzip reader; is this a gzip file?", zap.Error(err))
	}
	defer gzr.Close()

	comms := make(chan *goncrete.Communication)
	// goroutine that loads from the file
	go func(q chan<- *goncrete.Communication) {
		var added int
		buf := thrift.NewTMemoryBuffer()
		compact := thrift.NewTCompactProtocol(buf)
		trdr := tar.NewReader(gzr)
		for {
			_, localErr := trdr.Next()
			if localErr == io.EOF {
				break
			} else if localErr != nil {
				log.Error("error reading tar file", zap.Error(localErr))
			}

			io.Copy(buf, trdr)
			comm := goncrete.NewCommunication()
			if localErr = comm.Read(compact); localErr != nil {
				log.Error("error reading comm", zap.Error(localErr))
			}

			q <- comm
			added++
			buf.Reset()
		}

		close(q)
		log.Info("finished loading file", zap.Int("added-comms", added))
	}(comms)

	sock, err := thrift.NewTSocket(*fetchServer)
	if err != nil {
		log.Fatal("error opening socket to server", zap.Error(err))
	}

	transport := transFactory.GetTransport(sock)
	defer transport.Close()
	if err = transport.Open(); err != nil {
		log.Fatal("failed open transport", zap.Error(err))
	}

	storeCli := goncrete.NewStoreCommunicationServiceClientFactory(transport, protoFactory)
	alive, err := storeCli.Alive()
	if err != nil {
		log.Fatal("failed connect store server", zap.Error(err))
	}
	if !alive {
		log.Fatal("store not alive")
	}

	var ctr, errCtr int
	logPlusCtr := log.With(zap.Int("n-comms", ctr), zap.Int("n-errs", errCtr))
	for item := range comms {
		if err = storeCli.Store(item); err != nil {
			logPlusCtr.Error("error storing comm", zap.String("id", item.ID), zap.Error(err))
			errCtr++
		} else {
			ctr++
			if ctr%1000 == 0 {
				logPlusCtr.Info("ingesting")
			}
		}
	}
	logPlusCtr.Info("done")
}
