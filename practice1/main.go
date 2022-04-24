package main

import (
	"fmt"
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
	fmt.Printf("Start\n")
	eurodiffusion.CalculateEuroDiffusionForTestCase(testCases[0])

}
