package postgres

import (
	"fmt"
	"log"
	"testing"
)

func Test_NewConn(t *testing.T) {
	conn, err := NewConn()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(conn.Ping())
}

const insertSourceQuery = `
INSERT INTO sources (name) VALUES ($1) RETURNING id
`

const insertCampaignsQuery = `
INSERT INTO campaigns (name) VALUES ($1) RETURNING id
`

const insertCampaignsSourcesQuery = `
INSERT INTO campaigns_sources (campaign_id, source_id) VALUES ($1, $2)
`

func Test_fill(t *testing.T) {
	conn, err := NewConn()
	if err != nil {
		t.Fatal(err)
	}

	for i := range 100 {
		row := conn.QueryRow(insertSourceQuery, fmt.Sprintf("source_%d", i))
		if row.Err() != nil {
			log.Println(row.Err())
			continue
		}

		var sid uint64
		if err = row.Scan(&sid); err != nil {
			log.Println(err)
			continue
		}

		for j := range 5 {
			r := conn.QueryRow(insertCampaignsQuery, fmt.Sprintf("campaign_%d_%d", sid, j))
			if err != nil {
				log.Println(err)
				continue
			}

			var cid uint64
			if err = r.Scan(&cid); err != nil {
				log.Println(err)
				continue
			}

			_, err = conn.Exec(insertCampaignsSourcesQuery, cid, sid)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
