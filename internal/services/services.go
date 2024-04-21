package services

import (
	"adtelligent-internship/internal/models"
	"adtelligent-internship/internal/storage"
	"log"
	"math/rand/v2"
)

type Services interface {
	Generate() error
	MostDemandedSources(limit int) ([]storage.DemandedSource, error)
	GetUnlinkedCampaigns() ([]models.Campaign, error)
	GetEntitiesNames() ([]string, error)
}

type services struct {
	storage storage.Storage
}

func New(storage storage.Storage) Services {
	return &services{
		storage: storage,
	}
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
