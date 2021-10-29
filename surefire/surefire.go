package surefire

type TestSuite struct {
	Name     string `xml:"name,attr"`
	Time     string `xml:"time,attr"`
	Tests    string `xml:"tests,attr"`
	Errors   string `xml:"errors,attr"`
	Skipped  string `xml:"skipped,attr"`
	Failures string `xml:"failures,attr"`
	Properties Properties `xml:"properties"`
	Cases []TestCase `xml:"testcase"`
}

type Properties struct {
	Properties []Property `xml:"property"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type TestCase struct {
	Name      string `xml:"name,attr"`
	ClassName string `xml:"classname,attr"`
	Time      string `xml:"time,attr"`
}
