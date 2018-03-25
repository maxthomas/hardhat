package store

import (
	"fmt"

	"github.com/hltcoe/goncrete"
)

// PrintStorer prints IDs of incoming communications
type PrintStorer struct {
}

var (
	// PrintStorer impls StoreCommunicationService
	_ goncrete.StoreCommunicationService = (*PrintStorer)(nil)
)

func (p *PrintStorer) Alive() (bool, error) {
	return true, nil
}

func (p *PrintStorer) About() (*goncrete.ServiceInfo, error) {
	si := goncrete.NewServiceInfo()
	si.Name = "PrintStorer"
	si.Version = "0.0.1"
	return si, nil
}

func (p *PrintStorer) Store(comm *goncrete.Communication) error {
	fmt.Printf("%v", comm.ID)
	return nil
}
