# TC:   https://github.com/SUSE/caasp-test/features/skuba_upgrade/skuba_upgrade_bsc#1167320.feature
# PR:   https://github.com/SUSE/skuba/pull/1018
# BSC:  https://bugzilla.suse.com/show_bug.cgi?id=1167320
# FEATURE: skuba addon upgrade

# You are expected to run this test on a cluster bootstrapped with kubernetes-1.15.2
Feature: Check if cluster upgrade is fine

Scenario: Checking if cluster exists
    Given "skuba" exist in gopath
    And VARIABLE "imba-cluster" I get from CONFIG
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

    When I run "cat skubaconf.yaml"
    And I replace Gangway Version in OUTPUT and save it into skubaconf.yaml file
    And I run "kubectl apply -f skubaconf.yaml"
    And I run "skuba addon upgrade plan" in VAR:"imba-cluster" directory
    Then the output contains "cilium" and "gangway"
    Then the output contains "->" and "1.5.1" and "2.1.0-rev4"
