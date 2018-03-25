package annotator

import (
	"errors"
	"time"

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
	_ goncrete.AnnotateCommunicationService = (*Copier)(nil)
)

func (c *Copier) Annotate(original *goncrete.Communication) (*goncrete.Communication, error) {
	c.log.Debug("called", zap.String("method", "Annotate"))
	if original == nil {
		return nil, errors.New("can't copy a nil communication")
	}
	c.log.Debug("communication", zap.String("id", original.ID))
	return original, nil
}

func (c *Copier) GetMetadata() (*goncrete.AnnotationMetadata, error) {
	c.log.Debug("called", zap.String("method", "GetMetadata"))
	amd := goncrete.NewAnnotationMetadata()
	amd.Tool = "copier"
	amd.Timestamp = time.Now().Unix()
	amd.KBest = 0
	return amd, nil
}

func (c *Copier) GetDocumentation() (string, error) {
	c.log.Debug("called", zap.String("method", "GetDocumentation"))
	return "Copies the original communication and returns", nil
}

func (c *Copier) Shutdown() (err error) {
	c.log.Debug("called", zap.String("method", "Shutdown"))
	return nil
}
