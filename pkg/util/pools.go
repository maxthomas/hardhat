package util

import (
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// DeserializerPool returns a sync.Pool of deserializers
func DeserializerPool(fact thrift.TProtocolFactory) sync.Pool {
	return sync.Pool{New: func() interface{} {
		dser := thrift.NewTDeserializer()
		mcompact := fact.GetProtocol(dser.Transport)
		dser.Protocol = mcompact
		return dser
	},
	}
}

// SerializerPool returns a sync.Pool of serializers
func SerializerPool(fact thrift.TProtocolFactory) sync.Pool {
	return sync.Pool{New: func() interface{} {
		ser := thrift.NewTSerializer()
		mcompact := fact.GetProtocol(ser.Transport)
		ser.Protocol = mcompact
		return ser
	},
	}
}
