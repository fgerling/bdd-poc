## Quick Start
```
go get github.com/cucumber/godog/cmd/godog
git clone https://github.com/fgerling/bdd-poc.git
cd bdd-poc

## quick project layout.
## main_test.go as starting point.
tree

## check hello world, and run it.
## see how we get hints where the implementation is written.
cat features/hello/hello_world.feature
godog features/hello/hello_world.feature

## run "real" test. Look how it fails.
cat features/cluster/status.feature
godog features/cluster
## get real cluster (prepared beforehand) and run again
mv ~/cluster .
ls cluster
godog features/cluster
## yeah. Green test!

## let us use the beauty of "BDD" and write a new scenario:
cat << EOF >> features/makefile.feature
    Scenario: Check that skuba version is untagged from git
      Given I have "skuba" in PATH
      When I run "skuba"
      Then the output contains "UNTAGGED"
EOF

## check which line number has out new scenario
cat -n features/makefile.feature

## run our new scenario (without adding go code)
## the line number at the end of the file should match our new scenario
godog features/makefile.feature:36

## I am skipping the report generation here.
## see https://w3.suse.de/~fgerling/report.html as sample

## comparing running tests with godog vs compiled "testrunner"
## important: use not the term "testrunner". It could be confusing
godog -o tester
ls -l tester
time ./tester features/hello
time godog features/hello


## DEMO END
``` 
