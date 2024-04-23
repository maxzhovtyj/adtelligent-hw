package services

import (
	"github.com/maxzhovtyj/adtelligent-hw/internal/models"
	"github.com/maxzhovtyj/adtelligent-hw/internal/storage"
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

type Services interface {
	Generate() error
	MostDemandedSources(limit int) ([]storage.DemandedSource, error)
	GetUnlinkedCampaigns() ([]models.Campaign, error)
	GetEntitiesNames() ([]string, error)
	GetSourceCampaigns(req *GetSourceCampaignsRequest) ([]models.Campaign, error)
}

type services struct {
	storage storage.Storage
	cache   CampaignsToSourceCache
}

func New(storage storage.Storage) Services {
	s := &services{
		storage: storage,
	}

	start := make(chan bool)

	go s.initCache(start)

	res := <-start

	log.Println("cache initialized")

	if !res {
		panic("can't init cache")
	}

	return s
}

func (s *services) initCache(startSig chan<- bool) {
	all, err := s.storage.GetAllSourceCampaigns()
	if err != nil {
		log.Printf("failed to refresh sources campaigns: %v\n", err)
		startSig <- false
	}

	s.cache.Refresh(all)

	startSig <- true

	for {
		start := time.Now()

		all, err = s.storage.GetAllSourceCampaigns()
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

var p sync.Pool

func Acquire() *GetSourceCampaignsRequest {
	r := p.Get()
	if r == nil {
		return new(GetSourceCampaignsRequest)
	}

	return r.(*GetSourceCampaignsRequest)
}

func Release(r *GetSourceCampaignsRequest) {
	r.Reset()
	p.Put(r)
}

type GetSourceCampaignsRequest struct {
	ID      int
	Domains []string
}

func (r *GetSourceCampaignsRequest) Reset() {
	r.ID = 0
	r.Domains = r.Domains[:0]
}

func (s *services) GetSourceCampaigns(req *GetSourceCampaignsRequest) ([]models.Campaign, error) {
	//campaigns, err := s.storage.GetSourceCampaigns(req.ID)
	//if err != nil {
	//	return nil, err
	//}

	//use caching, main tuning
	campaigns := s.cache.Get(req.ID)

	var filteredCampaigns []models.Campaign

	for _, c := range campaigns {
		hasAll := true

		for _, d := range req.Domains {
			if !c.DomainWhitelist.Has(d) {
				hasAll = false
				break
			}
		}

		if !hasAll {
			continue
		}

		filteredCampaigns = append(filteredCampaigns, c)
	}

	return filteredCampaigns, nil
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
