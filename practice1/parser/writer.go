package parser

import (
	"bytes"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/akrava/software-engineering-methodology-practices/practice1/common"
)

func WriteResultsOfTestCase(w io.Writer, testCase common.TestCaseResults, caseNumber int) error {
	buf := bytes.NewBufferString("Case number ")
	buf.WriteString(strconv.Itoa(caseNumber))
	sort.Slice(testCase, func(i, j int) bool {
		diff := testCase[i].NumberDays - testCase[j].NumberDays
		if diff == 0 {
			return strings.Compare(testCase[i].CountryName, testCase[j].CountryName) < 0
		} else {
			return diff < 0
		}
	})
	for _, country := range testCase {
		buf.WriteString("\n")
		buf.WriteString(country.CountryName)
		buf.WriteString(" ")
		buf.WriteString(strconv.Itoa(country.NumberDays))
	}
	buf.WriteString("\n")
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}
