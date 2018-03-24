package util

import (
	"sync"
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

	errs := make(chan error, 4)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ser := pool.Get().(*thrift.TSerializer)
			defer pool.Put(ser)
			ser.Transport.Reset()

			for item := range comms {
				bytez, err := ser.Write(item)
				if !assert.NoError(t, err) {
					errs <- err
					return
				}

				if !assert.NotNil(t, bytez) {
					errs <- err
					return
				}
			}
		}()
	}
	wg.Wait()
	if len(errs) > 0 {
		firstErr := <-errs
		assert.FailNow(t, "at least one error", firstErr.Error())
	}
}
