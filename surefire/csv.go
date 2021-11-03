package surefire

import (
	"encoding/xml"
	"io"
	"path/filepath"
)

func Header() []string {
	return []string{
		"module",
		"class",
		"test",
		"test duration [seconds]",
		"test suite duration [seconds]",
		"test suite tests [number]",
		"test suite errors [number]",
		"test suite skipped [number]",
		"test suite failures [number]",
		"basedir",
	}
}

func Records(r io.Reader) ([][]string, error) {
	var suite TestSuite
	err := xml.NewDecoder(r).Decode(&suite)
	if err != nil {
		return nil, err
	}

	var basedir string
	var module string
	for _, v := range suite.Properties.Properties {
		if v.Name == "basedir" {
			basedir = v.Value
			module = filepath.Base(v.Value)
		}
	}

	var records [][]string
	for _, c := range suite.Cases {
		records = append(records, []string{
			module,
			c.ClassName,
			c.Name,
			c.Time,
			suite.Time,
			suite.Tests,
			suite.Errors,
			suite.Skipped,
			suite.Failures,
			basedir,
		})
	}

	return records, nil
}
