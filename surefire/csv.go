package surefire

import (
	"encoding/xml"
	"io"
	"path/filepath"
)

func Header() []string {
	return []string{"module", "class", "test", "duration", "basedir"}
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
	for _, v := range suite.Cases {
		records = append(records, []string{module, v.ClassName, v.Name, v.Time, basedir})
	}

	return records, nil
}
