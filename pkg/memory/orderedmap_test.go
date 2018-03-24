package memory

import (
	"testing"

	"github.com/hltcoe/goncrete"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	comm1 := goncrete.NewCommunication()
	comm1.ID = "comm1"
	comm2 := goncrete.NewCommunication()
	comm2.ID = "comm2"

	mp := NewOrderedMap()
	mp.Add(comm1, comm2)
	// mp.Add(comm2)

	items := mp.Iterate()
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "comm1", items[0].ID)
	assert.Equal(t, "comm2", items[1].ID)

	items = mp.Iterate()
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "comm1", items[0].ID)
	assert.Equal(t, "comm2", items[1].ID)

	assert.Equal(t, 2, mp.Len())
	item, present := mp.Get("comm1")
	assert.True(t, present)
	assert.Equal(t, "comm1", item.ID)

	item, present = mp.Get("quix")
	assert.False(t, present)
}
