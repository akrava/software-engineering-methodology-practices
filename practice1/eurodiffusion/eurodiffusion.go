package eurodiffusion

import (
	"fmt"
	"math"

	"github.com/akrava/software-engineering-methodology-practices/practice1/common"
)

func CalculateEuroDiffusionForTestCase(countries common.CountryList) (int, error) {
	minXl, minYl, maxYh, maxXh := getNumberRowsCols(countries)
	countRows := maxYh - minYl + 1
	countCols := maxXh - minXl + 1
	citiesGrid := make([][]*common.City, countRows)
	for y := range citiesGrid {
		citiesGrid[y] = make([]*common.City, countCols)
	}
	offsetX := minXl
	offsetY := minYl
	for _, country := range countries {
		country.InitCities()
		for _, city := range country.Cities {
			citiesGrid[city.Y-offsetY][city.X-offsetX] = city
		}
	}

	fmt.Printf("%#v\n", citiesGrid)
	fmt.Printf("%d %d\n", countRows, countCols)
	fmt.Printf("%d %d\n", minXl, minYl)
	fmt.Printf("%d %d\n", maxXh, maxYh)

	// for y := countRows - 1; y >= 0; y-- {
	// 	fmt.Printf("|")
	// 	for x := 0; x < countCols; x++ {
	// 		city := citiesGrid[y][x]
	// 		if city != nil {
	// 			fmt.Printf(city.Country.Name[0:1])
	// 		} else {
	// 			fmt.Printf(" ")
	// 		}
	// 	}
	// 	fmt.Printf("|\n")
	// }

	countDays := 1

	for {
		sumsToTransferOnTheBegginingOfTheDay := make([][]map[*common.Country]int, countRows)
		for y := range sumsToTransferOnTheBegginingOfTheDay {
			sumsToTransferOnTheBegginingOfTheDay[y] = make([]map[*common.Country]int, countCols)
		}
		for y := 0; y < countRows; y++ {
			for x := 0; x < countCols; x++ {
				city := citiesGrid[y][x]
				if city != nil {
					sumsToTransfer := map[*common.Country]int{}
					for country, amount := range city.CoinsAmount {
						sumsToTransfer[country] = amount / common.RepresentativePortionPerCoin
					}
					sumsToTransferOnTheBegginingOfTheDay[y][x] = sumsToTransfer
				}
			}
		}

		for y := 0; y < countRows; y++ {
			for x := 0; x < countCols; x++ {
				city := citiesGrid[y][x]
				if city != nil {

					sumsToTransferCurCity := sumsToTransferOnTheBegginingOfTheDay[y][x]

					if x > 0 {
						cityWest := citiesGrid[y][x-1]
						transferMoney(city, cityWest, sumsToTransferCurCity)
					}
					if y > 0 {
						citySouth := citiesGrid[y-1][x]
						transferMoney(city, citySouth, sumsToTransferCurCity)
					}
					if x < countCols-1 {
						cityEast := citiesGrid[y][x+1]
						transferMoney(city, cityEast, sumsToTransferCurCity)
					}
					if y < countRows-1 {
						cityNorth := citiesGrid[y+1][x]
						transferMoney(city, cityNorth, sumsToTransferCurCity)
					}
				}
			}
		}

		isFullAllCountries := true
		for _, country := range countries {
			if country.IsCompleated {
				continue
			}
			if country.IsFull(countries) && !country.IsCompleated {
				country.IsCompleated = true
				country.NumberDaysToBeCompleted = countDays
			} else {
				isFullAllCountries = false
			}
		}
		if isFullAllCountries {
			break
		}
		countDays++
	}
	
	// citiesGrid1 := make([][]*common.City, maxYh + 1)
	// for y := range citiesGrid1 {
	// 	citiesGrid1[y] = make([]*common.City, maxXh + 1)
	// }
	// for _, country := range countries {
	// 	for _, city := range country.Cities {
	// 		citiesGrid1[city.Y][city.X] = city
	// 	}
	// }

	// for y := maxYh; y >= 0; y-- {
	// 	fmt.Printf("|")
	// 	for x := 0; x <= maxXh; x++ {
	// 		city := citiesGrid1[y][x]
	// 		if city != nil {
	// 			fmt.Printf(city.Country.Name[0:1])
	// 		} else {
	// 			fmt.Printf(" ")
	// 		}
	// 	}
	// 	fmt.Printf("|\n")
	// }

	return 0, nil
}

func transferMoney(cityFrom, cityTo *common.City, amounts map[*common.Country]int) {
	if cityTo != nil {
		for country, _ := range cityFrom.CoinsAmount {
			sumToTransfer := amounts[country]
			if sumToTransfer > 0 {
				cityFrom.CoinsAmount[country] -= sumToTransfer
				cityTo.CoinsAmount[country] += sumToTransfer
			}
		}
	}
}

func getNumberRowsCols(countries common.CountryList) (int, int, int, int) {
	minXl := math.MaxInt
	minYl := math.MaxInt
	maxYh := math.MinInt
	maxXh := math.MinInt
	for _, country := range countries {
		if country.Xl < minXl {
			minXl = country.Xl
		}
		if country.Yl < minYl {
			minYl = country.Yl
		}
		if country.Xh > maxXh {
			maxXh = country.Xh
		}
		if country.Yh > maxYh {
			maxYh = country.Yh
		}
	}

	return minXl, minYl, maxYh, maxXh
}
