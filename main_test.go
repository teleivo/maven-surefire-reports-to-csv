package main

import (
	"bytes"
	"strings"
	"testing"
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
