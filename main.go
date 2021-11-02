package main

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/teleivo/surefire-reports-to-csv/surefire"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	src := flags.String("src", "", "Source directory containing Maven Surefire XML reports")
	dest := flags.String("dest", "", "Destination directory where CSV will be written to")
	debug := flags.Bool("debug", false, "Print debug information")
	err := flags.Parse(args[1:])
	if err != nil {
		return err
	}
	if *src == "" {
		return errors.New("src must be provided")
	}
	if *dest == "" {
		return errors.New("dest must be provided")
	}

	return csvConverter{from: *src, log: out, debug: *debug}.to(*dest)
}

type csvConverter struct {
	from  string
	log   io.Writer
	debug bool
}

func (cc csvConverter) to(dest string) error {
	s, err := os.Stat(dest)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dest, 0750)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	if s != nil && !s.IsDir() {
		return fmt.Errorf("dest path exists but is not a directory %q", dest)
	}

	// TODO collect errors in slice and report all of them
	//nolint:errcheck
	filepath.WalkDir(cc.from, func(path string, d fs.DirEntry, _ error) error {
		// TODO what to do on err?
		if d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".xml" {
			return nil
		}

		err = writeCSV(path, filepath.Join(dest, CSVFilename(path)))
		if err != nil {
			fmt.Fprintf(cc.log, "Failed to convert %q due to %s\n", path, err)
		} else {
			if cc.debug {
				fmt.Fprintf(cc.log, "Converted %q\n", path)
			}
		}

		return nil
	})

	return err
}

func CSVFilename(file string) string {
	fn := filepath.Base(file)
	return strings.TrimSuffix(fn, filepath.Ext(fn)) + ".csv"
}

func writeCSV(file string, dest string) error {
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer w.Close()

	return convert(r, w)
}

func convert(r io.Reader, w io.Writer) error {
	var suite surefire.TestSuite
	err := xml.NewDecoder(r).Decode(&suite)
	if err != nil {
		return err
	}

	var basedir string
	var module string
	for _, v := range suite.Properties.Properties {
		if v.Name == "basedir" {
			basedir = v.Value
			module = filepath.Base(v.Value)
		}
	}

	// maven module, basedir, class name, test name, duration
	records := [][]string{
		{"module", "class", "test", "duration", "basedir"},
	}
	for _, v := range suite.Cases {
		records = append(records, []string{module, v.ClassName, v.Name, v.Time, basedir})
	}

	c := csv.NewWriter(w)
	for _, record := range records {
		if err := c.Write(record); err != nil {
			return err
		}
	}
	// Write any buffered data to the underlying writer (standard output).
	c.Flush()
	if err := c.Error(); err != nil {
		return err
	}

	return nil
}
