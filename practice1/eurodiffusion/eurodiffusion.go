package eurodiffusion

import (
	"math"

	"github.com/akrava/software-engineering-methodology-practices/practice1/common"
)

func CalculateEuroDiffusionForTestCase(countries common.CountryList) common.TestCaseResults {
	minXl, minYl, maxYh, maxXh := getMinXlMinYlMaxYhMaxMh(countries)
	countRows := maxYh - minYl + 1
	countCols := maxXh - minXl + 1
	// initialize trimmed grid of cities
	citiesGrid := make([][]*common.City, countRows)
	for y := range citiesGrid {
		citiesGrid[y] = make([]*common.City, countCols)
	}
	// put cities into grid
	for _, country := range countries {
		country.InitCities()
		for _, city := range country.Cities {
			citiesGrid[city.Y-minYl][city.X-minXl] = city
		}
	}
	// start euro diffusion
	countDays := 0
	for {
		// check if all countries are full and set completed flag and number of days to those ones, which are full now
		if allCompleted := checkIfAllCountriesAreFullAndSetCompletedOnesIfNeeded(countries, countDays); allCompleted {
			break
		}
		// increase number of days
		countDays++
		// get amounts of money to transfer for each city in all currencies at the beggining of the day
		amountsToTransfer := getAmountsToTransferAtTheBegginingOfTheDay(citiesGrid, countRows, countCols)
		// transfer money between cities
		for y := 0; y < countRows; y++ {
			for x := 0; x < countCols; x++ {
				city := citiesGrid[y][x]
				if city == nil {
					continue
				}
				amountsToTransferCurrentCity := amountsToTransfer[y][x]
				if x > 0 {
					cityWest := citiesGrid[y][x-1]
					transferAmountOfMoneyBetweenCities(city, cityWest, amountsToTransferCurrentCity)
				}
				if y > 0 {
					citySouth := citiesGrid[y-1][x]
					transferAmountOfMoneyBetweenCities(city, citySouth, amountsToTransferCurrentCity)
				}
				if x < countCols-1 {
					cityEast := citiesGrid[y][x+1]
					transferAmountOfMoneyBetweenCities(city, cityEast, amountsToTransferCurrentCity)
				}
				if y < countRows-1 {
					cityNorth := citiesGrid[y+1][x]
					transferAmountOfMoneyBetweenCities(city, cityNorth, amountsToTransferCurrentCity)
				}
			}
		}
	}
	return getResultsForTestCase(countries)
}

func getResultsForTestCase(countries common.CountryList) common.TestCaseResults {
	testCaseResults := common.TestCaseResults{}
	for _, country := range countries {
		countryEuroDiffusion := common.CountryEuroDiffusion{
			CountryName: country.Name,
			NumberDays: country.NumberDaysToBeCompleted,
		}
		testCaseResults = append(testCaseResults, countryEuroDiffusion)
	}
	return testCaseResults
}

func checkIfAllCountriesAreFullAndSetCompletedOnesIfNeeded(countries common.CountryList, currentDay int) bool {
	allCountriesAreFull := true
	for _, country := range countries {
		if !country.IsCompleted && country.IsFull(countries) {
			country.IsCompleted = true
			country.NumberDaysToBeCompleted = currentDay
		}
		if !country.IsCompleted {
			allCountriesAreFull = false
		}
	}
	return allCountriesAreFull
}

func getAmountsToTransferAtTheBegginingOfTheDay(cities [][]*common.City, rows, cols int) [][]map[*common.Country]int {
	amountsToTransferOnTheBegginingOfTheDay := make([][]map[*common.Country]int, rows)
	for y := range amountsToTransferOnTheBegginingOfTheDay {
		amountsToTransferOnTheBegginingOfTheDay[y] = make([]map[*common.Country]int, cols)
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if city := cities[y][x]; city != nil {
				amountToTransferForCity := map[*common.Country]int{}
				for country, amount := range city.CoinsAmount {
					amountToTransferForCity[country] = amount / common.RepresentativePortionPerCoin
				}
				amountsToTransferOnTheBegginingOfTheDay[y][x] = amountToTransferForCity
			}
		}
	}
	return amountsToTransferOnTheBegginingOfTheDay
}

func transferAmountOfMoneyBetweenCities(cityFrom, cityTo *common.City, amounts map[*common.Country]int) {
	if cityTo == nil {
		return
	}
	for country := range cityFrom.CoinsAmount {
		sumToTransfer := amounts[country]
		if sumToTransfer > 0 {
			cityFrom.CoinsAmount[country] -= sumToTransfer
			cityTo.CoinsAmount[country] += sumToTransfer
		}
	}
}

func getMinXlMinYlMaxYhMaxMh(countries common.CountryList) (int, int, int, int) {
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
