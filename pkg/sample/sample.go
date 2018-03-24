package sample

import (
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
)

const (
	sampleUUID = "123e4567-e89b-12d3-a456-426655440000"
)

// Communication returns a sample communication
func Communication() *goncrete.Communication {
	comm := goncrete.NewCommunication()
	comm.ID = "sample"
	comm.UUID = UUID()
	comm.Metadata = Metadata()
	comm.Text = thrift.StringPtr("sample communication text")
	comm.Type = "document"
	return comm
}

// Metadata returns a sample valid AnnotationMetadata
func Metadata() *goncrete.AnnotationMetadata {
	md := goncrete.NewAnnotationMetadata()
	md.Tool = "sample"
	md.Timestamp = time.Now().Unix()
	return md
}

// UUID returns a sample goncrete.UUID
func UUID() *goncrete.UUID {
	uuid := goncrete.NewUUID()
	uuid.UuidString = sampleUUID
	return uuid
}
