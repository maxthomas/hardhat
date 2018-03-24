package thrift

import (
	"strconv"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// HostPort returns a string suitable for use in thrift.NewTSocket
func HostPort(host string, port uint64) string {
	return host + ":" + strconv.FormatUint(port, 10)
}

// NewSocket uses the specified host and port
// and returns a thrift.TSocket pointer, error tuple
func NewSocket(host string, port uint64) (*thrift.TSocket, error) {
	return thrift.NewTSocket(HostPort(host, port))
}

// NewServerSocket uses the specified host and port
// and returns a thrift.TServerSocket pointer, error tuple
func NewServerSocket(host string, port uint64) (*thrift.TServerSocket, error) {
	return thrift.NewTServerSocket(HostPort(host, port))
}
