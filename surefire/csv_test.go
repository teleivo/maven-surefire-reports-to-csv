package surefire

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestConvert(t *testing.T) {
	tc := map[string]struct {
		input string
		want  [][]string
		err   bool
	}{
		"ReportWithSkippedTestCase": {
			input: `<?xml version="1.0" encoding="UTF-8"?>
<testsuite xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://maven.apache.org/surefire/maven-surefire-plugin/xsd/surefire-test-report-3.0.xsd" version="3.0" name="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0.003" tests="1" errors="0" skipped="1" failures="0">
  <properties>
    <property name="java.specification.version" value="11"/>
    <property name="file.separator" value="/"/>
    <property name="basedir" value="/home/runner/work/dhis2-core/dhis2-core/dhis-2/dhis-services/dhis-service-administration"/>
  </properties>
  <testcase name="" classname="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0">
    <skipped/>
  </testcase>
</testsuite>`,
			want: [][]string{
				{
					"dhis-service-administration",
					"org.hisp.dhis.maintenance.HardDeleteAuditTest",
					"",
					"0",
					"0.003",
					"1",
					"0",
					"1",
					"0",
					"/home/runner/work/dhis2-core/dhis2-core/dhis-2/dhis-services/dhis-service-administration",
				},
			},
		},
		"ReportWithMultipleTests": {
			input: `<?xml version="1.0" encoding="UTF-8"?>
<testsuite xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://maven.apache.org/surefire/maven-surefire-plugin/xsd/surefire-test-report-3.0.xsd" version="3.0" name="org.hisp.dhis.analytics.data.AnalyticsServiceTest" time="171.217" tests="4" errors="1" skipped="0" failures="1">
  <properties>
    <property name="java.specification.version" value="11"/>
    <property name="maven.wagon.http.pool" value="false"/>
    <property name="java.vendor.url" value="http://www.azul.com/"/>
    <property name="user.timezone" value="Etc/UTC"/>
    <property name="java.vm.specification.version" value="11"/>
    <property name="os.name" value="Linux"/>
    <property name="user.home" value="/home/runner"/>
    <property name="user.language" value="en"/>
    <property name="file.separator" value="/"/>
    <property name="basedir" value="/home/runner/work/dhis2-core/dhis2-core/dhis-2/dhis-services/dhis-service-analytics"/>
  </properties>
  <testcase name="testMappingAggregation" classname="org.hisp.dhis.analytics.data.AnalyticsServiceTest" time="46.089">
    <system-out><![CDATA[* INFO  15:32:37,771 Found 7 analytics table types: [ORG_UNIT_TARGET, EVENT, VALIDATION_RESULT, DATA_VALUE, COMPLETENESS_TARGET, COMPLETENESS, ENROLLMENT] (DefaultAnalyticsTableGenerator.java [main])
* INFO  15:32:37,772 Analytics table update: AnalyticsTableUpdateParams{last years=null, skip resource tables=false, skip table types=[], skip programs=[], start time=2021-10-29T15:32:37} (DefaultAnalyticsTableGenerator.java [main])
* INFO  15:32:37,772 Last successful analytics table update: 'null' (DefaultAnalyticsTableGenerator.java [main])
* INFO  15:32:37,778 Generating resource table: '_orgunitstructure' (JdbcResourceTableStore.java [main])
* INFO  15:32:37,808 Resource table '_orgunitstructure' update done: '00:00:00.030' (JdbcResourceTableStore.java [main])
* INFO  15:33:21,705 Analytics tables dropped (DefaultAnalyticsTableService.java [main])
]]></system-out>
  </testcase>
  <testcase name="queryValidationResultTable" classname="org.hisp.dhis.analytics.data.AnalyticsServiceTest" time="41.134"/>
  <testcase name="testGridAggregation" classname="org.hisp.dhis.analytics.data.AnalyticsServiceTest" time="42.103"/>
  <testcase name="testSetAggregation" classname="org.hisp.dhis.analytics.data.AnalyticsServiceTest" time="41.879"/>
</testsuite>`,
			want: [][]string{
				{
					"dhis-service-analytics",
					"org.hisp.dhis.analytics.data.AnalyticsServiceTest",
					"testMappingAggregation",
					"46.089",
					"171.217",
					"4",
					"1",
					"0",
					"1",
					"/home/runner/work/dhis2-core/dhis2-core/dhis-2/dhis-services/dhis-service-analytics",
				},
				{
					"dhis-service-analytics",
					"org.hisp.dhis.analytics.data.AnalyticsServiceTest",
					"queryValidationResultTable",
					"41.134",
					"171.217",
					"4",
					"1",
					"0",
					"1",
					"/home/runner/work/dhis2-core/dhis2-core/dhis-2/dhis-services/dhis-service-analytics",
				},
				{
					"dhis-service-analytics",
					"org.hisp.dhis.analytics.data.AnalyticsServiceTest",
					"testGridAggregation",
					"42.103",
					"171.217",
					"4",
					"1",
					"0",
					"1",
					"/home/runner/work/dhis2-core/dhis2-core/dhis-2/dhis-services/dhis-service-analytics",
				},
				{
					"dhis-service-analytics",
					"org.hisp.dhis.analytics.data.AnalyticsServiceTest",
					"testSetAggregation",
					"41.879",
					"171.217",
					"4",
					"1",
					"0",
					"1",
					"/home/runner/work/dhis2-core/dhis2-core/dhis-2/dhis-services/dhis-service-analytics",
				},
			},
		},
		"ReportWithoutProperties": {
			input: `<?xml version="1.0" encoding="UTF-8"?>
<testsuite xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://maven.apache.org/surefire/maven-surefire-plugin/xsd/surefire-test-report-3.0.xsd" version="3.0" name="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0.003" tests="1" errors="0" skipped="1" failures="0">
  <testcase name="" classname="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0">
    <skipped/>
  </testcase>
</testsuite>`,
			want: [][]string{
				{
					"",
					"org.hisp.dhis.maintenance.HardDeleteAuditTest",
					"",
					"0",
					"0.003",
					"1",
					"0",
					"1",
					"0",
					"",
				},
			},
		},
		"ReportWithoutPropertyBaseDir": {
			input: `<?xml version="1.0" encoding="UTF-8"?>
<testsuite xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://maven.apache.org/surefire/maven-surefire-plugin/xsd/surefire-test-report-3.0.xsd" version="3.0" name="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0.003" tests="1" errors="0" skipped="1" failures="0">
  <testcase name="" classname="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0">
  <properties>
    <property name="java.specification.version" value="11"/>
    <property name="file.separator" value="/"/>
  </properties>
    <skipped/>
  </testcase>
</testsuite>`,
			want: [][]string{
				{
					"",
					"org.hisp.dhis.maintenance.HardDeleteAuditTest",
					"",
					"0",
					"0.003",
					"1",
					"0",
					"1",
					"0",
					"",
				},
			},
		},
		"ReportWithEmptyPropertyBaseDir": {
			input: `<?xml version="1.0" encoding="UTF-8"?>
<testsuite xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://maven.apache.org/surefire/maven-surefire-plugin/xsd/surefire-test-report-3.0.xsd" version="3.0" name="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0.003" tests="1" errors="0" skipped="1" failures="0">
  <testcase name="" classname="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0">
  <properties>
    <property name="java.specification.version" value="11"/>
    <property name="file.separator" value="/"/>
    <property name="basedir" value=""/>
  </properties>
    <skipped/>
  </testcase>
</testsuite>`,
			want: [][]string{
				{
					"",
					"org.hisp.dhis.maintenance.HardDeleteAuditTest",
					"",
					"0",
					"0.003",
					"1",
					"0",
					"1",
					"0",
					"",
				},
			},
		},
		"InvalidXML": {
			input: `<?xml version="1.0" encoding="UTF-8"?>
<testsuite xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://maven.apache.org/surefire/maven-surefire-plugin/xsd/surefire-test-report-3.0.xsd" version="3.0" name="org.hisp.dhis.maintenance.HardDeleteAuditTest" time="0.003" tests="1" errors="0" skipped="1" failures="0">
</tes>`,
			want: nil,
			err:  true,
		},
	}

	for k, v := range tc {
		t.Run(k, func(t *testing.T) {
			r := strings.NewReader(v.input)

			got, err := records(r)
			if v.err && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !v.err && err != nil {
				t.Fatalf("expected no error but got: %s", err)
			}

			if diff := cmp.Diff(v.want, got); diff != "" {
				t.Errorf("convert() mismatch (-want +got): \n%s", diff)
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
			c := CsvConverter{From: v.input, Concat: v.concat, Log: &w}

			dest := t.TempDir()

			err := c.To(dest)
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
		c := CsvConverter{From: "testdata/missing_src_directory/", Concat: false, Log: &w}

		err := c.To(t.TempDir())

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
		c := CsvConverter{From: src, Concat: false, Log: &w}

		err = c.To(dest)
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
		c := CsvConverter{From: "testdata/input", Concat: false, Log: &w}

		err = c.To(dest)

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
		c := CsvConverter{From: "testdata/input", Concat: false, Log: &w}

		err = c.To(dest)
		if err == nil {
			t.Error("expected an error but got none")
		}
		_, err = os.Stat(dest)
		if err == nil {
			// if destiation does not exist os.Stat errs, if parent dir has exec
			// bit not set it errs
			t.Fatal("destination should not be created. expected an error but got none")
		}
	})
}

func TestConcatConverter(t *testing.T) {
	t.Run("FailsIfToCannotBeCreated", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "dest")
		err := os.Mkdir(dest, 0440)
		if err != nil {
			t.Fatalf("failed to create dest dir for test: %s", err)
		}

		cc := &concatConverter{to: dest, once: &sync.Once{}}

		err = cc.convert("testdata/input/TEST-org.hisp.dhis.maintenance.HardDeleteAuditTest.xml")
		if err == nil {
			t.Error("expected an error but got none")
		}
		_, err = os.Stat(filepath.Join(dest, "surefire.csv"))
		if err == nil {
			// if destiation does not exist os.Stat errs, if parent dir has exec
			// bit not set it errs
			t.Fatal("surefire.csv should not be created. expected an error but got none")
		}
	})

	t.Run("FailsIfFromCannotBeRead", func(t *testing.T) {
		src := filepath.Join(t.TempDir(), "src")
		err := os.Mkdir(src, 0750)
		if err != nil {
			t.Fatalf("failed to create src dir for test: %s", err)
		}
		f := filepath.Join(src, "non-readable")
		err = os.WriteFile(f, []byte("data"), 0004)
		if err != nil {
			t.Fatalf("failed to create non-readable file for test: %s", err)
		}

		cc := &concatConverter{to: t.TempDir(), once: &sync.Once{}}

		err = cc.convert(f)
		if err == nil {
			t.Error("expected an error but got none")
		}
		// the to file will currently still be created and not removed if the
		// only file to convert fails conversion
		_, err = os.Stat(filepath.Join(cc.to, "surefire.csv"))
		if err != nil {
			t.Fatalf("surefire.csv should have been created. expected no error but got %s", err)
		}
	})
}
