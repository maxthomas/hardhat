package summarizer

import (
	"errors"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
	"go.uber.org/zap"
)

// Copier implements goncrete.AnnotateCommunicationService by returning
// a copy of the original communication, or an error if nil.
type Copier struct {
	log *zap.Logger
}

// NewCopier returns an instantiated Copier
func NewCopier(log *zap.Logger) *Copier {
	return &Copier{log}
}

var (
	_ goncrete.SummarizationService = (*Copier)(nil)
)

func (c *Copier) Summarize(query *goncrete.SummarizationRequest) (*goncrete.Summary, error) {
	if query == nil {
		return nil, errors.New("can't summarize nil request")
	}
	if !query.IsSetSourceCommunication() {
		return nil, errors.New("SourceCommuncation required to be set")
	}
	summ := goncrete.NewSummary()
	summ.SummaryCommunication = query.GetSourceCommunication()
	return summ, nil
}

func (c *Copier) GetCapabilities() ([]*goncrete.SummarizationCapability, error) {
	caps := goncrete.NewSummarizationCapability()
	caps.Type = goncrete.SummarySourceType_DOCUMENT
	caps.Lang = "all"
	items := []*goncrete.SummarizationCapability{caps}
	return items, nil
}

// About returns information about the service
func (c *Copier) About() (*goncrete.ServiceInfo, error) {
	si := goncrete.NewServiceInfo()
	si.Name = "copier"
	si.Version = "latest"
	si.Description = thrift.StringPtr("Returns a summary with a copied communication")
	return si, nil
}

// Alive answers: is the service alive?
func (c *Copier) Alive() (bool, error) {
	return true, nil
}
