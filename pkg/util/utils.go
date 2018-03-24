package util

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
)

var (
	protocolFactory = thrift.NewTCompactProtocolFactory()
)

// LoadTarFile takes in a tar.Reader and pushes the files (not dirs)
// to a channel
func LoadTarFile(tarRdr *tar.Reader, data chan<- []byte) error {
	var err error
	for {
		_, err = tarRdr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error reading tar file: %v", err)
		}

		dataBytes, err := ioutil.ReadAll(tarRdr)
		if err != nil {
			return err
		}
		data <- dataBytes
	}
	return nil
}

// LoadCommunication takes in a byte array and returns a Communication or error
func LoadCommunication(c []byte) (*goncrete.Communication, error) {
	memBuf := thrift.NewTMemoryBuffer()
	deser := thrift.NewTDeserializer()
	deser.Protocol = protocolFactory.GetProtocol(memBuf)
	deser.Transport = memBuf
	comm := goncrete.NewCommunication()
	err := deser.Read(comm, c)
	return comm, err
}
