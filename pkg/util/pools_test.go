package util

import (
	"testing"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
	"github.com/maxthomas/hardhat/pkg/sample"
	"github.com/stretchr/testify/assert"
)

var (
	fact = thrift.NewTCompactProtocolFactory()
)

func TestSerPool(t *testing.T) {
	pool := SerializerPool(fact)

	nItems := 100
	comms := make(chan *goncrete.Communication, nItems)
	for i := 0; i < nItems; i++ {
		comms <- sample.Communication()
	}
	close(comms)
	for i := 0; i < 3; i++ {
		go func() {
			for item := range comms {
				ser := pool.Get().(*thrift.TSerializer)

				bytez, err := ser.Write(item)
				if assert.NoError(t, err) {
					assert.FailNow(t, "failed parallel serialization")
				}

				assert.NotNil(t, bytez)
				ser.Transport.Reset()
				pool.Put(ser)
			}
		}()
	}
}
