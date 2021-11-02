# sure (Maven Surefire Reports To CSV)

[![Build and Test](https://github.com/teleivo/maven-surefire-reports-to-csv/actions/workflows/build_test.yml/badge.svg)](https://github.com/teleivo/maven-surefire-reports-to-csv/actions/workflows/build_test.yml)
[![golangci-lint](https://github.com/teleivo/maven-surefire-reports-to-csv/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/teleivo/maven-surefire-reports-to-csv/actions/workflows/golangci-lint.yml)
[![codecov](https://codecov.io/gh/teleivo/maven-surefire-reports-to-csv/branch/main/graph/badge.svg?token=1VFP7UVS4Z)](https://codecov.io/gh/teleivo/maven-surefire-reports-to-csv)
[![Release](https://img.shields.io/github/release/teleivo/maven-surefire-reports-to-csv.svg)](https://github.com/teleivo/maven-surefire-reports-to-csv/releases/latest)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg)](https://github.com/goreleaser)

`sure` helps analyze test count and duration in Java projects built using
[Maven](https://maven.apache.org/) and
[Maven Surefire Plugin](https://maven.apache.org/surefire/maven-surefire-plugin/).

Maven Surefire Plugin generates reports after running tests in a Maven module.
Surefire can be configured to write the reports in XML. `sure` allows you
to convert them into CSV. This way you can easily analyze the number of tests,
their duration either individually or for example per Maven module in any tool
you whish.

## Usage

Download a binary for your platform at
[releases](https://github.com/teleivo/maven-surefire-reports-to-csv/releases)

And convert Surefire XML reports to CSV

```sh
sure \
  -src ~/code/yourproject \
  -dest ./here
```

### Compile

If you have [Go](https://golang.org/) installed and want to compile yourself
:smile: you can

Run it directly using

```sh
go run main.go \
  -src ~/code/yourproject \
  -dest ./here
```

Or build a binary first

```sh
go build -o sure
```

## Limitation

* Was only tested on Maven Surefire reports of schema
  "https://maven.apache.org/surefire/maven-surefire-plugin/xsd/surefire-test-report-3.0.xsd"
* You need to concatenate the CSV files into a single file yourself if that is
  what you prefer for your analysis. Be aware that each CSV has a header that
  you will need to ignore. This might work for you :smile;

    cat ~/somewhere_nice/* | grep -v module > combined.csv
