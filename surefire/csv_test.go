package surefire

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

			got, err := Records(r)
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
