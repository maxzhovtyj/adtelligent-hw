package storage

import (
	"context"
	"fmt"
	"github.com/maxzhovtyj/adtelligent-hw/internal/models"
	"time"
)

const randomSourcesQuery = `
INSERT INTO sources (name)
SELECT concat('Source #', md5(random()::text))
FROM generate_series(1, $1)
RETURNING id, name;
`

func (s *storage) RandomSources(count int) ([]models.Source, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	rows, err := s.db.QueryContext(ctx, randomSourcesQuery, count)
	if err != nil {
		return nil, fmt.Errorf("failed to insert random sources: %w", err)
	}

	var sources []models.Source

	for rows.Next() {
		var src models.Source

		if err = rows.Scan(&src.ID, &src.Name); err != nil {
			return nil, err
		}

		sources = append(sources, src)
	}

	return sources, nil
}

const selectMostDemandedSourcesQuery = `
SELECT source_id, count(source_id) AS count
FROM campaigns_sources
GROUP BY source_id
ORDER BY count DESC
LIMIT $1;
`

type DemandedSource struct {
	ID    string
	Count int
}

func (s *storage) GetMostDemandedSources(limit int) ([]DemandedSource, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	rows, err := s.db.QueryContext(ctx, selectMostDemandedSourcesQuery, limit)
	if err != nil {
		return nil, err
	}

	var mostDemanded []DemandedSource

	for rows.Next() {
		var ds DemandedSource

		if err = rows.Scan(&ds.ID, &ds.Count); err != nil {
			return nil, err
		}

		mostDemanded = append(mostDemanded, ds)
	}

	return mostDemanded, nil
}

const selectSourceCampaigns = `
SELECT c.id, c.name 
FROM sources s
INNER JOIN campaigns_sources cs ON cs.source_id = s.id
INNER JOIN campaigns c ON cs.campaign_id = c.id
WHERE s.id = $1
`

func (s *storage) GetSourceCampaigns(sourceID int) ([]models.Campaign, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	rows, err := s.db.QueryContext(ctx, selectSourceCampaigns, sourceID)
	if err != nil {
		return nil, err
	}

	var campaigns []models.Campaign

	for rows.Next() {
		var cmp models.Campaign

		if err = rows.Scan(&cmp.ID, &cmp.Name); err != nil {
			return nil, err
		}

		campaigns = append(campaigns, cmp)
	}

	return campaigns, nil
}

const selectAllSourceCampaigns = `
SELECT s.id, c.id, c.name 
FROM sources s
INNER JOIN campaigns_sources cs ON cs.source_id = s.id
INNER JOIN campaigns c ON cs.campaign_id = c.id
`

func (s *storage) GetAllSourceCampaigns() (map[int][]models.Campaign, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	rows, err := s.db.QueryContext(ctx, selectAllSourceCampaigns)
	if err != nil {
		return nil, err
	}

	all := make(map[int][]models.Campaign)

	for rows.Next() {
		var sourceID, campaignID int
		var campaignName string

		if err = rows.Scan(&sourceID, &campaignID, &campaignName); err != nil {
			return nil, err
		}

		all[sourceID] = append(all[sourceID], models.Campaign{
			ID:   campaignID,
			Name: campaignName,
		})
	}

	return all, nil
}
