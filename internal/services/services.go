package services

import (
	"github.com/maxzhovtyj/adtelligent-hw/internal/models"
	"github.com/maxzhovtyj/adtelligent-hw/internal/storage"
	"log"
	"math/rand/v2"
	"time"
)

type Services interface {
	Generate() error
	MostDemandedSources(limit int) ([]storage.DemandedSource, error)
	GetUnlinkedCampaigns() ([]models.Campaign, error)
	GetEntitiesNames() ([]string, error)
	GetSourceCampaigns(sourceId ...int) ([]models.Campaign, error)
}

type services struct {
	storage storage.Storage
	cache   CampaignsToSourceCache
}

func New(storage storage.Storage) Services {
	s := &services{
		storage: storage,
	}

	go s.initCache()

	return s
}

func (s *services) initCache() {
	for {
		start := time.Now()

		all, err := s.storage.GetAllSourceCampaigns()
		if err != nil {
			log.Printf("failed to refresh sources campaigns: %v\n", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		s.cache.Refresh(all)
		log.Printf("refreshed %d sources campaigns in %s\n", len(all), time.Since(start))
		time.Sleep(1 * time.Minute)
	}
}

func (s *services) GetSourceCampaigns(sourceIds ...int) ([]models.Campaign, error) {
	if len(sourceIds) == 0 {
		return nil, nil
	}

	return s.storage.GetSourceCampaigns(sourceIds[0])

	// use caching, main tuning
	//return s.cache.Get(sourceIds[0]), nil
}

const (
	// each source should have [0;10] linked campaigns
	sourceMaxLinkedCampaigns = 10
)

// Generate task 1.0
func (s *services) Generate() error {
	randomSources, err := s.storage.RandomSources(100)
	if err != nil {
		return err
	}

	randomCampaigns, err := s.storage.RandomCampaigns(100)
	if err != nil {
		return err
	}

	for _, src := range randomSources {
		connectedCampaigns := rand.N(sourceMaxLinkedCampaigns + 1)

		for range connectedCampaigns {
			randomCmp := randomCampaigns[rand.N(len(randomCampaigns))]

			err = s.storage.CreateCampaignToSourceLink(randomCmp.ID, src.ID)
			if err != nil {
				log.Printf("failed to link campaign #d to source #d: %v\n", err)
				continue
			}
		}
	}

	return nil
}

// MostDemandedSources task 1.1
func (s *services) MostDemandedSources(limit int) ([]storage.DemandedSource, error) {
	return s.storage.GetMostDemandedSources(limit)
}

// GetUnlinkedCampaigns task 1.2
func (s *services) GetUnlinkedCampaigns() ([]models.Campaign, error) {
	return s.storage.GetUnlinkedCampaigns()
}

// GetEntitiesNames task 1.3
func (s *services) GetEntitiesNames() ([]string, error) {
	return s.storage.GetEntitiesNames()
}
