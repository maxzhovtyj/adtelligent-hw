package storage

import (
	"adtelligent-internship/internal/models"
	"context"
	"fmt"
	"time"
)

const randomCampaignsQuery = `
INSERT INTO campaigns (name)
SELECT concat('Campaign #', md5(random()::text))
FROM generate_series(1, $1)
RETURNING id, name
`

func (s *storage) RandomCampaigns(count int) ([]models.Campaign, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	rows, err := s.db.QueryContext(ctx, randomCampaignsQuery, count)
	if err != nil {
		return nil, fmt.Errorf("failed to insert random campaigns: %w", err)
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

const insertCampaignToSourceQuery = `
INSERT INTO campaigns_sources (campaign_id, source_id) VALUES ($1, $2)
`

func (s *storage) CreateCampaignToSourceLink(campaignID, sourceID int) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	_, err := s.db.ExecContext(ctx, insertCampaignToSourceQuery, campaignID, sourceID)
	if err != nil {
		return err
	}

	return nil
}

const selectUnlinkedCampaignsQuery = `
SELECT c.id, c.name
FROM campaigns c
         LEFT JOIN campaigns_sources cs ON c.id = cs.campaign_id
WHERE cs.id IS NULL
`

func (s *storage) GetUnlinkedCampaigns() ([]models.Campaign, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	rows, err := s.db.QueryContext(ctx, selectUnlinkedCampaignsQuery)
	if err != nil {
		return nil, err
	}

	var unlinkedCampaigns []models.Campaign

	for rows.Next() {
		var cmp models.Campaign

		if err = rows.Scan(&cmp.ID, &cmp.Name); err != nil {
			return nil, err
		}

		unlinkedCampaigns = append(unlinkedCampaigns, cmp)
	}

	return unlinkedCampaigns, nil
}
