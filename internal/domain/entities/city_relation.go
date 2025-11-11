package entities

type CityRelation struct {
	CityFullName string
	ProfileID    int64
}

func NewCityRelation(cityFullName string, profileID int64) *CityRelation {
	return &CityRelation{
		CityFullName: cityFullName,
		ProfileID:    profileID,
	}
}
