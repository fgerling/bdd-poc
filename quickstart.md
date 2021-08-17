# Quick Start
I want to guide you through the first steps of BDD and godog.
We will prepare our workstation, checkout the project, execute the first BDD test and add a new scenario.

## Prerequise
As prerequise we need a working golang installation. Either follow the the [upstream documentation](https://golang.org/doc/install) or stay SUSE:
```
sudo zypper install go
eval "$(go env)"
PATH=$GOPATH/bin:$PATH
# you can also modify the PATH in your .bashrc
```

## Get the code and godog
With a working golang installation, we need some test cases and golang code that implement the test case steps.
```
git clone https://github.com/fgerling/bdd-poc.git
cd bdd-poc
```

Then, we need a framework that can match the "steps" of test cases against golang functions. We are using [`godog`](https://github.com/cucumber/godog).
```
go get github.com/cucumber/godog/cmd/godog
```

## Project layout
Have a look at the project layout:
- `features` is the place where all the test cases are stored.
- `main_test.go` is the entry point for godog, and is the place where the BDD steps are matched against golang functions.
- In `internal/*` the implementation of the BDD steps can be found.


## Execute the first BDD test case
Let us check our first "Hello World" test case and execute it.
```
cat features/hello/hello_world.feature
godog features/hello/hello_world.feature
```

## Write first BDD test case
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
