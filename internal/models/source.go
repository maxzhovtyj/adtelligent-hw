package models

type Source struct {
	ID        int        `json:"ID"`
	Name      string     `json:"name"`
	Campaigns []Campaign `json:"campaigns,omitempty"`
}
