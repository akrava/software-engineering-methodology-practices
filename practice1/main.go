package main

import (
	"log"
	"os"

	"github.com/akrava/software-engineering-methodology-practices/practice1/eurodiffusion"
	"github.com/akrava/software-engineering-methodology-practices/practice1/parser"
)

func main() {
	testCases, err := parser.ParseTestCasesFromInputStream(os.Stdin)
	if err != nil {
		log.Printf("Couldn't parse test cases from input: %v", err)
		os.Exit(1)
	}
	for i, testCase := range testCases {
		results := eurodiffusion.CalculateEuroDiffusionForTestCase(testCase)
		if err := parser.WriteResultsOfTestCase(os.Stdout, results, i+1); err != nil {
			log.Printf("Couldn't write results of test case #%d: %v", i+1, err)
			os.Exit(1)
		}
	}
}
