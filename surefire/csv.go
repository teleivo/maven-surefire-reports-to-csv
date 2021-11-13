package surefire

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type CsvConverter struct {
	From   string
	Concat bool
	Log    io.Writer
	Debug  bool
}

func (cc CsvConverter) To(dest string) error {
	s, err := os.Stat(dest)
	if err == nil && !s.IsDir() {
		return fmt.Errorf("dest path exists but is not a directory %q", dest)
	}
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dest, 0750)
	}
	if err != nil {
		return err
	}

	var converter converter
	if cc.Concat {
		converter = &concatConverter{to: dest, once: &sync.Once{}}
	} else {
		converter = &separateConverter{to: dest}
	}
	defer converter.Close()

	// TODO collect errors in slice and report all of them
	// using WalkDir as godoc of Walk declares it as being more efficient
	err = filepath.WalkDir(cc.From, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if cc.From == path {
				// stop the Walk if from cannot be read
				return fmt.Errorf("failed to walk %q: %w", cc.From, err)
			}
			fmt.Fprintf(cc.Log, "Failed to process %q due to %s\n", path, err)
			return nil
		}
		if d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".xml" {
			return nil
		}

		err = converter.convert(path)
		if err != nil {
			fmt.Fprintf(cc.Log, "Failed to convert %q due to %s\n", path, err)
			return nil
		}
		if cc.Debug {
			fmt.Fprintf(cc.Log, "Converted %q\n", path)
		}

		return nil
	})

	return err
}

type converter interface {
	convert(from string) error
	io.Closer
}

type concatConverter struct {
	to string
	// TODO
	w    io.WriteCloser
	csv  *csv.Writer
	once *sync.Once
}

type separateConverter struct {
	to string
}

func (cc *concatConverter) convert(from string) error {
	var err error
	cc.once.Do(func() {
		// declaration needed so err is closed over. w, err := ... does not work
		var w *os.File
		w, err = os.Create(path.Join(cc.to, "surefire.csv"))
		if err != nil {
			return
		}

		cc.w = w
		cc.csv = csv.NewWriter(w)
		err = cc.csv.Write(header())
	})
	if err != nil {
		return err
	}

	r, err := os.Open(from)
	if err != nil {
		return err
	}
	defer r.Close()

	rr, err := records(r)
	if err != nil {
		return err
	}
	for _, record := range rr {
		if err := cc.csv.Write(record); err != nil {
			return err
		}
	}

	cc.csv.Flush()
	if err := cc.csv.Error(); err != nil {
		return err
	}

	return nil
}

func (cc *concatConverter) Close() error {
	// TODO do I need to keep the io.Writer so I can close it?
	if cc.csv != nil {
		return cc.w.Close()
	}
	return nil
}

func (sc *separateConverter) convert(from string) error {
	r, err := os.Open(from)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(filepath.Join(sc.to, CSVFilename(from)))
	if err != nil {
		return err
	}
	defer w.Close()

	c := csv.NewWriter(w)
	if err := c.Write(header()); err != nil {
		return err
	}

	rr, err := records(r)
	if err != nil {
		return err
	}
	for _, record := range rr {
		if err := c.Write(record); err != nil {
			return err
		}
	}

	c.Flush()
	if err := c.Error(); err != nil {
		return err
	}

	return nil
}

func (sc *separateConverter) Close() error {
	return nil
}

func CSVFilename(file string) string {
	fn := filepath.Base(file)
	return strings.TrimSuffix(fn, filepath.Ext(fn)) + ".csv"
}

func header() []string {
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

func records(r io.Reader) ([][]string, error) {
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
