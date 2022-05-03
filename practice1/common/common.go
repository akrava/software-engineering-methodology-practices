package common

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
	IsCompleted            bool
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
	NumberDays int
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
