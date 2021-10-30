# sure (Maven Surefire Reports To CSV)

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
./sure \
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

* This project is more like a script with typesafety. I did not write any
  tests. I did only run sanity checks on the CSV. Use at your own risk.

* You need to concatenate the CSV files into a single file yourself if that is
  what you prefer for your analysis. Be aware that each CSV has a header that
  you will need to ignore. This might work for you :smile;

    cat ~/somewhere_nice/* | grep -v module > combined.csv
