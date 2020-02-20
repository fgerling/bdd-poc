# bdd-poc
## Setup

```sh
go get github.com/cucumber/godog/cmd/godog
git clone https://github.com/fgerling/bdd-poc.git
cd bdd-poc
```

## Test
```sh
# test whole feature
$ godog makefile.feature
# [...]

# test specific scenario
$ godog makefile.feature:11
Feature: Skuba repository

  Scenario: skuba make test                   # makefile.feature:11
    Given the "skuba" repository exsist       # makefile_test.go:49 -> theRepositoryExsist
    When I run "make test" in "skuba"         # makefile_test.go:31 -> iRunIn
    Then the file "skuba/coverage.out" exsist # makefile_test.go:41 -> theFileExsist

1 scenarios (1 passed)
3 steps (3 passed)
39.599326885s
```
