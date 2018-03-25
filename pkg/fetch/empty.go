package fetch

import (
	"errors"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hltcoe/goncrete"
	"go.uber.org/zap"
)

type EmptyFetcher struct {
	logger *zap.Logger
}

func NewEmptyFetcher(logger *zap.Logger) *EmptyFetcher {
	return &EmptyFetcher{logger}
}

var (
	_       goncrete.FetchCommunicationService = (*EmptyFetcher)(nil)
	commIDs                                    = []string{"1", "2", "3"}
)

func (g *EmptyFetcher) Alive() (bool, error) {
	g.logger.Info("called", zap.String("method", "Alive"))
	return true, nil
}

func (g *EmptyFetcher) About() (*goncrete.ServiceInfo, error) {
	g.logger.Info("called", zap.String("method", "About"))
	si := goncrete.NewServiceInfo()
	si.Name = "EmptyFetcher"
	si.Version = "0.0.1"
	return si, nil
}

func metadata() *goncrete.AnnotationMetadata {
	md := goncrete.NewAnnotationMetadata()
	md.Tool = "EmptyFetcher"
	md.Timestamp = time.Now().Unix()
	return md
}

func (g *EmptyFetcher) Fetch(req *goncrete.FetchRequest) (*goncrete.FetchResult_, error) {
	g.logger.Info("called", zap.String("method", "Fetch"))
	if req == nil {
		return nil, errors.New("can't fetch nil request")
	}

	fr := goncrete.NewFetchResult_()
	fr.Communications = make([]*goncrete.Communication, 3)
	for _, id := range req.CommunicationIds {
		comm := goncrete.NewCommunication()
		uid := goncrete.NewUUID()
		uid.UuidString = "123e4567-e89b-12d3-a456-426655440000"
		comm.UUID = uid
		comm.Type = "document"
		comm.Metadata = metadata()
		comm.Text = thrift.StringPtr("empty")

		switch id {
		case "1":
			comm.ID = "1"
		case "2":
			comm.ID = "2"
		case "3":
			comm.ID = "3"
		default:
			continue
		}

		fr.Communications = append(fr.Communications, comm)
	}
	return fr, nil
}

func (g *EmptyFetcher) GetCommunicationCount() (int64, error) {
	g.logger.Debug("called", zap.String("method", "GetCommunicationCount"))
	return int64(3), nil
}

func (g *EmptyFetcher) GetCommunicationIDs(off int64, count int64) ([]string, error) {
	g.logger.Debug("called",
		zap.String("method", "GetCommunicationIDs"),
		zap.Int64("offset", off), zap.Int64("count", count),
	)
	return commIDs, nil
}
