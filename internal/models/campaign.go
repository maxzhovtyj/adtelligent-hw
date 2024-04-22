package models

type Campaign struct {
	ID              int
	Name            string
	DomainWhitelist *Whitelist
}

func (c *Campaign) InitWhitelist(values ...string) {
	c.DomainWhitelist = &Whitelist{
		data: make(map[string]struct{}),
	}

	c.DomainWhitelist.Put(values...)
}
