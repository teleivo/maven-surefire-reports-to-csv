package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/teleivo/surefire-reports-to-csv/surefire"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	// ExitOnError makes the error message look cleaner to the user
	// but makes testing hard. ContinueOnError allows me to capture the
	// returned error. Unfortunately, flag will print the error and usage and
	// main() will print the error again.
	// TODO is there a way to handle this better?
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	src := flags.String("src", "", "Source directory containing Maven Surefire XML reports.")
	dest := flags.String("dest", "", "Destination directory where CSV(s) will be written to. It will be created if does not exist.")
	concat := flags.Bool("concat", false, "Concatenate all Maven Surefire XML reports into one CSV file.")
	debug := flags.Bool("debug", false, "Print debug information.")
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

	return csvConverter{
		from:   *src,
		concat: *concat,
		log:    out,
		debug:  *debug,
	}.to(*dest)
}

type csvConverter struct {
	from   string
	concat bool
	log    io.Writer
	debug  bool
}

func (cc csvConverter) to(dest string) error {
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
	if cc.concat {
		converter = &concatConverter{to: dest, once: &sync.Once{}}
	} else {
		converter = &separateConverter{to: dest}
	}
	defer converter.Close()

	// TODO collect errors in slice and report all of them
	// using WalkDir as godoc of Walk declares it as being more efficient
	err = filepath.WalkDir(cc.from, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if cc.from == path {
				// stop the Walk if from cannot be read
				return fmt.Errorf("failed to walk %q: %w", cc.from, err)
			}
			fmt.Fprintf(cc.log, "Failed to process %q due to %s\n", path, err)
			return nil
		}
		if d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".xml" {
			return nil
		}

		err = converter.convert(path)
		if err != nil {
			fmt.Fprintf(cc.log, "Failed to convert %q due to %s\n", path, err)
			return nil
		}
		if cc.debug {
			fmt.Fprintf(cc.log, "Converted %q\n", path)
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
	// TODO error handling how to?
	var err error
	cc.once.Do(func() {
		w, err := os.Create(path.Join(cc.to, "surefire.csv"))
		if err != nil {
			return
		}

		cc.w = w
		cc.csv = csv.NewWriter(w)
		err = cc.csv.Write(surefire.Header())
	})
	if err != nil {
		return err
	}

	r, err := os.Open(from)
	if err != nil {
		return err
	}
	defer r.Close()

	rr, err := surefire.Records(r)
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
	if err := c.Write(surefire.Header()); err != nil {
		return err
	}

	rr, err := surefire.Records(r)
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
