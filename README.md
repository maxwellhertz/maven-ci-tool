A toy tool to automate some CI steps for Maven projects. Currently, this tool does the following things:

1. Update the POM to include maven-surefire-plugin and jacoco-maven-plugin if they don't already exist;
2. Execute `mvn clean test` to run all the unit tests and generate a report about code coverage.