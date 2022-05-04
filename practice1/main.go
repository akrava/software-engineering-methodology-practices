package main

import (
	"log"
	"os"

	"github.com/akrava/software-engineering-methodology-practices/practice1/common"
	"github.com/akrava/software-engineering-methodology-practices/practice1/eurodiffusion"
	"github.com/akrava/software-engineering-methodology-practices/practice1/parser"
)

func main() {
	testCases, err := parser.ParseTestCasesFromInputStream(os.Stdin)
	if err != nil {
		log.Printf("Couldn't parse test cases from input: %v", err)
		os.Exit(1)
	}
	testCasesResults := []common.TestCaseResults{}
	for i, testCase := range testCases {
		if err := common.VerifyAllCountriesAreNeighborhoodsAndNonOverlapping(testCase); err != nil {
			log.Printf("Test case #%d is invalid: %v", i, err)
			os.Exit(1)
		}
		results := eurodiffusion.CalculateEuroDiffusionForTestCase(testCase)
		testCasesResults = append(testCasesResults, results)
	}
	for i, testCaseResult := range testCasesResults {
		if err := parser.WriteResultsOfTestCase(os.Stdout, testCaseResult, i+1); err != nil {
			log.Printf("Couldn't write results of test case #%d: %v", i+1, err)
			os.Exit(1)
		}
	}
}
