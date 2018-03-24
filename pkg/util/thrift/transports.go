package thrift

import "git.apache.org/thrift.git/lib/go/thrift"

func TransportFactory(buffered bool, framed bool) thrift.TTransportFactory {
	var transportFactory thrift.TTransportFactory
	if buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}
	return transportFactory
}
