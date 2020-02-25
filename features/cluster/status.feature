# doc: https://github.com/SUSE/skuba/blob/master/README.md

Feature: Skuba cluster status

  Scenario: checkout cluster status
    Given there is "cluster" directory
    And "skuba" exist in gopath
    When I run "skuba cluster status" in "cluster" directory
    Then the output contains "master" and "worker"
