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

	err = os.Mkdir(*dest, 0750)
	if err != nil {
		return err
	}

	// TODO collect errors in slice and report all of them
	//nolint:errcheck
	filepath.WalkDir(*src, func(path string, d fs.DirEntry, _ error) error {
		// TODO what to do on err?
		if d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".xml" {
			return nil
		}

		err = writeCSV(path, filepath.Join(*dest, CSVFilename(path)))
		if err != nil {
			fmt.Fprintf(out, "Failed to convert %q due to %s\n", path, err)
		} else {
			fmt.Fprintf(out, "Converted %q\n", path)
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
