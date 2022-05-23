package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akrava/software-engineering-methodology-practices/practice2/metrics"
)

func main() {
	if modulePath, err := getModulePath(); err != nil {
		log.Fatal(err)
	} else {
		moduleMetrics := metrics.ModuleMetrics{Path: modulePath}
		if err := moduleMetrics.CountAllMetrics(); err != nil {
			log.Fatal(err)
		}
		moduleMetrics.CalculateCommentSaturationLevel()
		fmt.Print(moduleMetrics.String())
	}
}

func getModulePath() (string, error) {
	if len(os.Args) != 2 {
		log.Printf("Pass the correct path to the library or module in arguments. Usage: %s <module path>", os.Args[0])
		return "", fmt.Errorf("expected 1 argument, got %d", len(os.Args)-1)
	}
	return os.Args[1], nil
}
