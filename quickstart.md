## Quick Start
I want to guide you through the first steps of BDD. 

First, we need a framework that can match "steps" of test cases against golang functions. We are using [`godog`](https://github.com/cucumber/godog).
```
go get github.com/cucumber/godog/cmd/godog
```

Next, we need some test cases and golang code that implement the test case steps.
```
git clone https://github.com/fgerling/bdd-poc.git
cd bdd-poc
```

Have a look at the project layout:
- `features` is the place where all the test cases are stored.
- `main_test.go` is the entry point for godog, and is the place where the BDD steps are matched against golang functions.
- In `internal/*` the implementation of the BDD steps can be found.


Let us check our first "Hello World" test case and execute it.
```
cat features/hello/hello_world.feature
godog features/hello/hello_world.feature
```

We can also write a new scenario based on the "steps" we have seen in the hello world example.
```
cat << EOF >> features/makefile.feature

    Scenario: Check that skuba version is untagged from git
      Given I have "skuba" in PATH
      When I run "skuba version"
      Then the output contains "UNTAGGED"
EOF

## check which line number has our new scenario
cat -n features/makefile.feature

## run our new scenario (without adding go code)
## the line number at the end of the file should match our new scenario
godog features/makefile.feature:36
```

The output of `godog` gives you a hint where the implementation of one of the BDD steps is located. 
Check for example `internal/suse/suse.go:16` for the implementation of "I have in PATH".

