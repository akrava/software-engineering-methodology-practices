package common

import "fmt"

const (
	MaxNumberOfCountries         = 20
	MaxCountryNameLength         = 25
	MinCoordValue                = 1
	MaxCoordValue                = 10
	InitialCoinsAmount           = 1_000_000
	RepresentativePortionPerCoin = 1_000
)

type Country struct {
	Name                    string
	Xl                      int
	Yl                      int
	Xh                      int
	Yh                      int
	Cities                  []*City
	IsCompleted             bool
	NumberDaysToBeCompleted int
}

type CountryList []*Country

type TestCases []CountryList

type City struct {
	X           int
	Y           int
	CoinsAmount map[*Country]int
	Country     *Country
}

type CountryEuroDiffusion struct {
	CountryName string
	NumberDays  int
}

type TestCaseResults []CountryEuroDiffusion

func (c *Country) InitCities() {
	c.Cities = []*City{}
	for x := c.Xl; x <= c.Xh; x++ {
		for y := c.Yl; y <= c.Yh; y++ {
			city := &City{
				X:           x,
				Y:           y,
				CoinsAmount: map[*Country]int{c: InitialCoinsAmount},
				Country:     c,
			}
			c.Cities = append(c.Cities, city)
		}
	}
}

func (c *Country) IsFull(countries CountryList) bool {
	for _, city := range c.Cities {
		if len(city.CoinsAmount) != len(countries) {
			return false
		} else {
			for country, amount := range city.CoinsAmount {
				isAnyCountry := false
				for _, c := range countries {
					if c == country {
						isAnyCountry = true
						break
					}
				}
				if !isAnyCountry {
					return false
				}
				if amount <= 0 {
					return false
				}
			}
		}
	}
	return true
}

func VerifyAllCountriesAreNeighborhoodsAndNonOverlapping(countries CountryList) error {
	for _, country := range countries {
		if country.Xl > country.Xh || country.Yl > country.Yh {
			return fmt.Errorf("country %s has got invalid coords", country.Name)
		}
		countryHasGotNeighborhoods := false
		for _, anotherCountry := range countries {
			if country == anotherCountry {
				continue
			}
			if country.Xh < anotherCountry.Xl || country.Xl > anotherCountry.Xh ||
				country.Yh < anotherCountry.Yl || country.Yl > anotherCountry.Yh {
				if !countryHasGotNeighborhoods {
					countryHasGotNeighborhoods =
						((anotherCountry.Xl-country.Xh == 1 || country.Xl-anotherCountry.Xh == 1) &&
							anotherCountry.Yl <= country.Yh && anotherCountry.Yh >= country.Yl) ||
							((anotherCountry.Yl-country.Yh == 1 || country.Yl-anotherCountry.Yh == 1) &&
								anotherCountry.Xl <= country.Xh && anotherCountry.Xh >= country.Xl)
				}
			} else {
				return fmt.Errorf("country %s is overlapping with country %s", country.Name, anotherCountry.Name)
			}
		}
		if !countryHasGotNeighborhoods && len(countries) > 1 {
			return fmt.Errorf("country %s doesn't have got neighborhoods", country.Name)
		}
	}
	return nil
}
