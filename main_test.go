package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRun(t *testing.T) {
	tc := map[string]struct {
		args []string
		err  string
	}{
		"UnkownFlagFailsParsing": {
			args: []string{
				"sure",
				"-unkownflag",
			},
			err: "-unkownflag",
		},
		"SrcIsMandatory": {
			args: []string{
				"sure",
				"-dest",
				t.TempDir(),
			},
			err: "src must be provided",
		},
		"DestIsMandatory": {
			args: []string{
				"sure",
				"-src",
				t.TempDir(),
			},
			err: "dest must be provided",
		},
	}

	for k, tc := range tc {
		t.Run(k, func(t *testing.T) {
			var out bytes.Buffer

			err := run(tc.args, &out)

			if tc.err != "" {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("got error %q but want %q", err, tc.err)
				}
			}
			if tc.err == "" && err != nil {
				t.Fatalf("expected no error but got: %s", err)
			}
		})
	}
}

func TestCsvConverter(t *testing.T) {
	tc := map[string]struct {
		input  string
		concat bool
		want   string
	}{
		"OneCSVFilePerXML": {
			input:  "testdata/input",
			concat: false,
			want:   "testdata/expected/separate",
		},
		"ConcatenatedCSVFile": {
			input:  "testdata/input",
			concat: true,
			want:   "testdata/expected/concat",
		},
	}

	for k, v := range tc {
		t.Run(k, func(t *testing.T) {
			var w bytes.Buffer
			c := csvConverter{from: v.input, concat: v.concat, log: &w}

			dest := t.TempDir()

			err := c.to(dest)
			if err != nil {
				t.Fatalf("expected no error but got %s", err)
			}

			de, err := os.ReadDir(dest)
			if err != nil {
				t.Fatalf("failed to read dest dir %q due to %s", dest, err)
			}

			want, err := os.ReadDir(v.want)
			if err != nil {
				t.Fatalf("failed to read dest dir %q due to %s", dest, err)
			}

			var wantNames []string
			for _, d := range want {
				wantNames = append(wantNames, d.Name())
			}
			var gotNames []string
			for _, got := range de {
				if !got.Type().IsRegular() {
					t.Errorf("expected regular file to be created instead got %v", got)
				}
				gotNames = append(gotNames, got.Name())
			}
			// ensure that the exact files are produces, no less, no more
			less := func(a, b string) bool { return a < b }
			if diff := cmp.Diff(wantNames, gotNames, cmpopts.SortSlices(less)); diff != "" {
				t.Errorf("to() file mismatch (-want +got): \n%s", diff)
			}

			// ensure that the content of the csv file is as epected
			for _, d := range want {
				want, err := os.ReadFile(path.Join(v.want, d.Name()))
				if err != nil {
					t.Errorf("expected no error but got %s", err)
					continue
				}
				got, err := os.ReadFile(path.Join(dest, d.Name()))
				if err != nil {
					t.Errorf("expected no error but got %s", err)
					continue
				}
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("to() mismatch (-want +got): \n%s", diff)
				}
			}
		})
	}

	t.Run("FailsIfSrcDoesNotExist", func(t *testing.T) {
		var w bytes.Buffer
		c := csvConverter{from: "testdata/missing_src_directory/", concat: false, log: &w}

		err := c.to(t.TempDir())

		if err == nil {
			t.Fatal("expected an error but got none")
		}
	})

	t.Run("IgnoresSubdirWhichCannotBeRead", func(t *testing.T) {
		src := t.TempDir()
		subdir := filepath.Join(src, "sealed")
		err := os.Mkdir(subdir, 7)
		if err != nil {
			t.Fatalf("failed to create root read-only %q due to %s", subdir, err)
		}
		b, err := ioutil.ReadFile("./testdata/input/TEST-org.hisp.dhis.maintenance.HardDeleteAuditTest.xml")
		if err != nil {
			t.Fatalf("failed to read report due to %s", err)
		}
		err = ioutil.WriteFile(filepath.Join(src, "TEST-org.hisp.dhis.maintenance.HardDeleteAuditTest.xml"), b, 0440)
		if err != nil {
			t.Fatalf("failed to write report due to %s", err)
		}

		dest := t.TempDir()

		var w bytes.Buffer
		c := csvConverter{from: src, concat: false, log: &w}

		err = c.to(dest)
		if err != nil {
			t.Errorf("expected no error but got: %s", err)
		}
		if !strings.HasPrefix(w.String(), "Failed to process") {
			t.Errorf("want log to start with 'Failed to process', instead got %q", w.String())
		}

		// ensure other report is still converted
		de, err := os.ReadDir(dest)
		if err != nil {
			t.Fatalf("failed to read dest dir %q due to %s", dest, err)
		}

		if got, want := len(de), 1; got != want {
			t.Fatalf("got %d files/dirs, want %d file", got, want)
		}

		got := de[0]
		if !got.Type().IsRegular() {
			t.Errorf("expected regular file to be created instead got %v", got)
		}
		if diff := cmp.Diff("TEST-org.hisp.dhis.maintenance.HardDeleteAuditTest.csv", got.Name()); diff != "" {
			t.Errorf("to() file mismatch (-want +got): \n%s", diff)
		}
	})

	t.Run("FailsIfDestExistsButIsNotADirectory", func(t *testing.T) {
		destDir := t.TempDir()
		dest := filepath.Join(destDir, "surefire")
		err := ioutil.WriteFile(dest, []byte("some content"), 0400)
		if err != nil {
			t.Fatalf("failed to create temp file for test: %s", err)
		}

		var w bytes.Buffer
		c := csvConverter{from: "testdata/input", concat: false, log: &w}

		err = c.to(dest)

		if err == nil {
			t.Fatal("expected an error but got none")
		}
		if want := "dest path exists but is not a directory"; !strings.Contains(err.Error(), want) {
			t.Fatalf("got %q but want it to contain %q", err, want)
		}
	})

	t.Run("FailsIfDestDoesNotExistAndCannotBeCreated", func(t *testing.T) {
		d := t.TempDir()
		parent := filepath.Join(d, "parent")
		err := os.Mkdir(parent, 0550)
		if err != nil {
			t.Fatalf("failed to create dest dir for test: %s", err)
		}
		dest := filepath.Join(parent, "dest")

		var w bytes.Buffer
		c := csvConverter{from: "testdata/input", concat: false, log: &w}

		err = c.to(dest)
		if err == nil {
			t.Fatal("expected an error but got none")
		}
		_, err = os.Stat(dest)
		if err == nil {
			// if destiation does not exist os.Stat errs, if parent dir has exec
			// bit not set it errs
			t.Fatal("destination should not be created. expected an error but got none")
		}
	})
}
