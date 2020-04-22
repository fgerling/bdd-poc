## Quick Start
I want to guide you through the first steps of BDD. 

First we need a framework that can match "Steps" of BDD test cases against golang functions. We choosed [`godog`](https://github.com/cucumber/godog).
```
go get github.com/cucumber/godog/cmd/godog
```
Next we get this repository with some sample test cases and golang code that tests those cases.
```
git clone https://github.com/fgerling/bdd-poc.git
cd bdd-poc
```

Have a look at the project layout:
- `features` is the place where all the test cases are stored.
- `main_test.go` is the entry point for godog, and is the place where the BDD steps are matched against golang functions.
- `internal/*`, the actual implementation of BDD steps, helper functions, etc. happens here. 


Now let us check our first "Hello World" test case, and run it.
```
cat features/hello/hello_world.feature
godog features/hello/hello_world.feature
```


Let us use the beauty of "BDD" and write a new scenario, and run it:
```
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
```
The output of `godog` gives you a hont where the implementation of one of the BDD steps is located. 
Check for example `internal/suse/suse.go:16` for the implementation of "Given I have in PATH".
