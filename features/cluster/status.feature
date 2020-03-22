# doc: https://github.com/SUSE/skuba/blob/master/README.md

Feature: Skuba cluster status

  Scenario: checkout cluster status
    Given VARIABLE "cluster" I get from CONFIG
    And "skuba" exist in gopath
    When I run "skuba cluster status" in VAR:"cluster" directory
    Then the output contains "master"
    And the output contains "worker"
