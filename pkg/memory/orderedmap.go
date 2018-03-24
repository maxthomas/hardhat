package memory

import (
	"sync"

	"github.com/hltcoe/goncrete"
)

type OrderedMap struct {
	mp    map[string]goncrete.Communication
	order []string

	mutex sync.RWMutex
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		mp:    make(map[string]goncrete.Communication),
		order: make([]string, 0),
	}
}

func (o *OrderedMap) Add(comms ...*goncrete.Communication) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	for _, comm := range comms {
		id := comm.GetID()
		o.mp[id] = *comm
		o.order = append(o.order, id)
	}
}

func (o *OrderedMap) Iterate() []*goncrete.Communication {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	sorted := make([]*goncrete.Communication, len(o.order))
	for i, id := range o.order {
		entry := o.mp[id]
		sorted[i] = &entry
	}
	return sorted
}

func (o *OrderedMap) Len() int {
	return len(o.mp)
}

func (o *OrderedMap) Get(id string) (goncrete.Communication, bool) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	item, present := o.mp[id]
	return item, present
}

func (o *OrderedMap) Slice(start, end int) []string {
	ids := make([]string, 0)
	currLen := len(o.order)
	if start < 0 || end <= 0 || start > currLen-1 {
		return ids
	}
	if end >= currLen {
		return o.order[start:]
	}
	return o.order[start:end]
}
