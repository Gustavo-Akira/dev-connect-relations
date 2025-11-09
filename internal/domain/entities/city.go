package entities

import "strings"

type City struct {
	Name    string
	Country string
	State   string
}

func NewCity(name, country, state string) *City {
	name = normalizeString(name)
	country = normalizeString(country)
	state = normalizeString(state)
	return &City{
		Name:    name,
		Country: country,
		State:   state,
	}
}

func normalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (c *City) GetFullName() string {
	return c.Name + ", " + c.State + ", " + c.Country
}
