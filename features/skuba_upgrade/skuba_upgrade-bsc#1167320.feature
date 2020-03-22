# TC:   https://github.com/fgerling/bdd-poc
# PR:   https://github.com/SUSE/skuba/pull/1018
# FEATURE: skuba addon upgrade

# You are expected to run this test on a cluster bootstrapped with kubernetes-1.15.2
Feature: Check if cluster upgrade is fine

Scenario: Checking if cluster exists
    Given "skuba" exist in gopath
    And VARIABLE "imba-cluster" equals "/Users/alexeitighineanu/go/src/github.com/fgerling/bdd-poc/imba-cluster"
    When I run "skuba cluster status" in VAR:"imba-cluster" directory
    Then the output contains "master" and "worker"
    Then the output contains "1.15.2"
    When I run "kubectl get all --namespace=kube-system"
    Then the output contains "cilium" and "dex"
    When I run "skuba version"
    Then the output contains "v1.2.6" or "v1.2.7"

    When I run "skuba addon upgrade plan" in VAR:"imba-cluster" directory
    Then the output contains "congratulations" or "are already"

    When I run "kubectl get configmaps skuba-config -n kube-system -o yaml" in VAR:"imba-cluster" directory
    And I replace Cilium Version in OUTPUT and save it into skubaconf.yaml file
    And I run "kubectl apply -f skubaconf.yaml"
    And I run "skuba addon upgrade plan" in VAR:"imba-cluster" directory
    Then the output contains "addon upgrades for" or "->"

    And I replace Gangway Version in OUTPUT and save it into skubaconf.yaml file
    And I run "kubectl apply -f skubaconf.yaml"
    And I run "skuba addon upgrade plan" in VAR:"imba-cluster" directory
    Then the output contains "addon upgrades for" and "gangway"