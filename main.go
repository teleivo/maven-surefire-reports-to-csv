package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

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

	return surefire.CsvConverter{
		From:   *src,
		Concat: *concat,
		Log:    out,
		Debug:  *debug,
	}.To(*dest)
}
