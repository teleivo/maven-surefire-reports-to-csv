package main

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
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

	// TODO could probably use a []err since there errors could accumulate but
	// I still want to convert as many reports into csv's as possible
	files, err := os.ReadDir(*src)
	for _, f := range files {
		if !f.IsDir() && strings.ToLower(filepath.Ext(f.Name())) != ".xml" {
			continue
		}

		err = writeCSV(filepath.Join(*src, f.Name()), filepath.Join(*dest, CSVFilename(f.Name())))
		if err != nil {
			fmt.Fprintf(out, "Failed to convert %q due to %s\n", f.Name(), err)
		} else {
			fmt.Fprintf(out, "Converted %q\n", f.Name())
		}
	}

	return err
}

func CSVFilename(file string) string {
	fn := filepath.Base(file)
	return strings.TrimSuffix(fn, filepath.Ext(fn)) + ".csv"
}

func writeCSV(file string, dest string) error {
	var suite surefire.TestSuite
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = xml.Unmarshal([]byte(data), &suite)
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

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			return err
		}
	}
	// Write any buffered data to the underlying writer (standard output).
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}
