package parser

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/akrava/software-engineering-methodology-practices/practice1/common"
)

type ReaderState int

const (
	NumberOfCountries ReaderState = iota
	CountryDefinition
)

func ParseTestCasesFromInputStream(r io.Reader) (common.TestCases, error) {
	scanner := bufio.NewScanner(r)
	testCases := common.TestCases{}
	curCountryList := common.CountryList{}
	curState := NumberOfCountries
	numberOfCountriesLeft := 0
	for scanner.Scan() {
		curLine := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading input: %v", err)
			return nil, err
		}
		switch curState {
		case NumberOfCountries:
			{
				if number, err := strconv.Atoi(strings.TrimSpace(curLine)); err != nil {
					log.Printf("Error parsing number of countries: `%v` is not a number", curLine)
					return nil, err
				} else if number < 0 || number > common.MaxNumberOfCountries {
					log.Printf("Number of countries `%v` should be: 1 ≤ c ≤ %d", curLine, common.MaxNumberOfCountries)
					return nil, errors.New("number of countries doesn't fit the allowed range")
				} else {
					numberOfCountriesLeft = number
					curState = CountryDefinition
				}
			}
		case CountryDefinition:
			{
				if country, err := parseCountryDefinition(curLine); err != nil {
					return nil, err
				} else {
					curCountryList = append(curCountryList, country)
				}
				numberOfCountriesLeft--
				if numberOfCountriesLeft == 0 {
					testCases = append(testCases, curCountryList)
					curCountryList = common.CountryList{}
					curState = NumberOfCountries
				}
			}
		}
		if curState == CountryDefinition && numberOfCountriesLeft == 0 {
			break
		}
	}
	return testCases, nil
}

func parseCountryDefinition(inputData string) (*common.Country, error) {
	countryTokens := strings.Fields(inputData)
	if len(countryTokens) != 5 {
		log.Printf("Error parsing country definition: expected 5 tokens in `%v`", inputData)
		return nil, errors.New("country definition requires name and 4 coords")
	}
	countryName := countryTokens[0]
	if len(countryName) > common.MaxCountryNameLength {
		log.Printf("Country name `%v` shouldn't exceed 25 characters", inputData)
		return nil, errors.New("country name is too long")
	}
	coords := [4]int{}
	for i := 1; i < 5; i++ {
		if number, err := strconv.Atoi(countryTokens[i]); err != nil {
			log.Printf("Error parsing coord: `%v` is not a number", countryTokens[i])
			return nil, err
		} else if number < common.MinCoordValue || number > common.MaxCoordValue {
			log.Printf("Coord %v should be: %d ≤ x ≤ %d", countryTokens[i], common.MinCoordValue, common.MaxCoordValue)
			return nil, errors.New("coordinate doesn't fit the allowed range")
		} else {
			coords[i-1] = number
		}
	}
	return &common.Country{
		Name: countryName,
		Xl:   coords[0],
		Yl:   coords[1],
		Xh:   coords[2],
		Yh:   coords[3],
	}, nil
}
