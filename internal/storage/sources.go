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
