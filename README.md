# Maven Surefire Reports To CSV

[Maven Surefire Plugin](https://maven.apache.org/surefire/maven-surefire-plugin/)
can create XML reports of Java test runs. This project allows you to convert
them into CSV. This way you can easily analyze the number of tests, their
duration individually or for example per Maven module.

## Limitation

* This project is more like a script with typesafety. I did not write any
  tests. I did only run sanity checks on the CSV. Use at your own risk.
* You need to copy the XML reports into a single directory if you want to
  analyze a multi-module project. This might work for you

    find . -name "TEST-*.xml" -exec cp {} ~/somewhere_nice \;

* You need to concatenate the CSV files into a single file yourself if that is
  what you prefer for your analysis. Be aware that each CSV has a header that
  you will need to ignore. This might work for you :smile;

    cat ~/somewhere_nice/* | grep -v module > combined.csv
