package main

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

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
}
