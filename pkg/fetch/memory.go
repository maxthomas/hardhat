package fetch

import (
	"github.com/hltcoe/goncrete"
	"github.com/maxthomas/hardhat/pkg/memory"
	"go.uber.org/zap"
)

type MemoryFetcher struct {
	om     *memory.OrderedMap
	logger *zap.Logger
}

func NewMemoryFetcher(om *memory.OrderedMap, logger *zap.Logger) *MemoryFetcher {
	return &MemoryFetcher{om, logger}
}

var (
	_ goncrete.FetchCommunicationService = (*MemoryFetcher)(nil)
)

func (g *MemoryFetcher) Alive() (bool, error) {
	g.logger.Info("called", zap.String("method", "alive"))
	return true, nil
}

func (g *MemoryFetcher) About() (*goncrete.ServiceInfo, error) {
	g.logger.Info("called", zap.String("method", "about"))
	si := goncrete.NewServiceInfo()
	si.Name = "MemoryFetcher"
	si.Version = "0.0.1"
	return si, nil
}

func (g *MemoryFetcher) Fetch(req *goncrete.FetchRequest) (*goncrete.FetchResult_, error) {
	g.logger.Info("called", zap.String("method", "fetch"))
	fr := goncrete.NewFetchResult_()
	fr.Communications = make([]*goncrete.Communication, 0)
	for _, id := range req.CommunicationIds {
		item, present := g.om.Get(id)
		if present {
			fr.Communications = append(fr.Communications, &item)
		}
	}
	return fr, nil
}

func (g *MemoryFetcher) GetCommunicationCount() (int64, error) {
	g.logger.Info("called", zap.String("method", "comms-count"))
	return int64(g.om.Len()), nil
}

func (g *MemoryFetcher) GetCommunicationIDs(off int64, count int64) ([]string, error) {
	g.logger.Info("called",
		zap.String("method", "get-comm-ids"),
		zap.Int64("offset", off), zap.Int64("count", count),
	)
	start, diff := int(off), int(count)
	ids := g.om.Slice(start, start+diff)
	g.logger.Debug("returning", zap.Strings("items", ids))
	return ids, nil
}
