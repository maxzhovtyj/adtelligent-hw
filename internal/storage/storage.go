package storage

import (
	"context"
	"database/sql"
	"github.com/maxzhovtyj/adtelligent-hw/internal/models"
	"time"
)

type Storage interface {
	RandomSources(count int) ([]models.Source, error)
	RandomCampaigns(count int) ([]models.Campaign, error)
	CreateCampaignToSourceLink(campaignID, sourceID int) error
	GetUnlinkedCampaigns() ([]models.Campaign, error)
	GetEntitiesNames() ([]string, error)
	GetMostDemandedSources(limit int) ([]DemandedSource, error)
}

type storage struct {
	db *sql.DB
}

const selectEntitiesNamesQuery = `
SELECT name FROM sources
UNION
SELECT name FROM campaigns
`

func (s *storage) GetEntitiesNames() ([]string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	rows, err := s.db.QueryContext(ctx, selectEntitiesNamesQuery)
	if err != nil {
		return nil, err
	}

	var all []string

	for rows.Next() {
		var name string

		if err = rows.Scan(&name); err != nil {
			return nil, err
		}

		all = append(all, name)
	}

	return all, nil
}

func New(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}
