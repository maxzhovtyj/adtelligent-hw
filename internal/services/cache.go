package services

import (
	"github.com/maxzhovtyj/adtelligent-hw/internal/models"
	"sync"
)

type CampaignsToSourceCache struct {
	cache map[int][]models.Campaign
	mx    sync.RWMutex
}

func (c *CampaignsToSourceCache) Get(id int) []models.Campaign {
	c.mx.RLock()
	campaigns := c.cache[id]
	c.mx.RUnlock()

	return campaigns[:]
}

func (c *CampaignsToSourceCache) Put(id int, campaign []models.Campaign) {
	c.mx.Lock()
	c.cache[id] = campaign
	c.mx.Unlock()
}

func (c *CampaignsToSourceCache) Refresh(n map[int][]models.Campaign) {
	c.mx.Lock()
	c.cache = n
	c.mx.Unlock()
}
