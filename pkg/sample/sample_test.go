package sample

import (
	"testing"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/stretchr/testify/assert"
)

var (
	fact = thrift.NewTCompactProtocolFactory()
)

func TestSerialization(t *testing.T) {
	ser := thrift.NewTSerializer()
	mcompact := thrift.NewTCompactProtocol(ser.Transport)
	ser.Protocol = mcompact

	structs := []thrift.TStruct{UUID(), Metadata(), Communication()}
	for _, item := range structs {
		bytez, err := ser.Write(item)
		assert.NoError(t, err)
		assert.NotNil(t, bytez)
	}
}
