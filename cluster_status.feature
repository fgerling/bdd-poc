# doc: https://github.com/SUSE/skuba/blob/master/README.md

Feature: Skuba cluster status

  Scenario: checkout cluster status
    Given there is "imba-cluster" directory
    And "skuba" exist in gopath
    When I run "skuba cluster status" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "master" and "worker"