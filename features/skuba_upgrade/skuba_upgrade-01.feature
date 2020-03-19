# TC:   https://github.com/fgerling/bdd-poc
# PR:   https://github.com/SUSE/skuba/pull/911
# FEATURE: skuba cluster upgrade

# You are expected to run this test on a cluster bootstrapped with kubernetes-1.15.2
Feature: Check if cluster upgrade is fine

Scenario: Checking if cluster exists
    Given "skuba" exist in gopath
    And VARIABLE "imba-cluster" equals "/root/go/src/github.com/fgerling/bdd-poc/imba-cluster"
    When I run "skuba cluster status" in VAR:"imba-cluster" directory
    Then the output contains "master" and "worker"
    Then the output contains "1.15.2"
    When I run "kubectl get all --namespace=kube-system"
    Then the output contains "cilium" and "dex"
    When I run "skuba version"
    Then the output contains "v1.2.6" or "v1.2.7"

    When I run "skuba cluster upgrade plan" in VAR:"imba-cluster" directory
    Then the output contains "current kubernetes" and "latest kubernetes"
    #Then the output contains "upgrade path to update"
    #Then the output contains "addon upgrades from"

    When I run "skuba addon upgrade apply"
    Then the output contains "congratulations"

Scenario: Applying upgrade on nodes
    When I run "kubectl get pods --namespace=kube-system"
    And VARIABLE "imba-cluster" equals "/root/go/src/github.com/fgerling/bdd-poc/imba-cluster"
    When VARIABLE "privileged-pods" equals ContainersFROMOutput "kured-"
    And VARIABLES "commandchecks" equals "kubectl describe pod -n kube-system " plus VAR:"privileged-pods"
    And I run VARS:"commandchecks" and IPSFromOutput

    #UPGRADING MASTERS FIRST
    When VARIABLES "commandupgrades" equals "skuba node upgrade plan " plus Master Nodes
    And I run UPGRADE VARS:"commandupgrades" in VAR:"imba-cluster" directory
    Then the output contains "apiserver" and "controller-manager" and "scheduler"
    And the output contains "etcd" and "kubelet" and "cri-o"
    And I run UPGRADE VARS:"commandupgrades" in VAR:"imba-cluster" directory
    Then the output contains "apiserver" and "controller-manager" and "scheduler"
    And the output contains "etcd" and "kubelet" and "cri-o"
    And I run UPGRADE VARS:"commandupgrades" in VAR:"imba-cluster" directory
    Then the output contains "apiserver" and "controller-manager" and "scheduler"
    And the output contains "etcd" and "kubelet" and "cri-o"

    When VARIABLES "upgradeapply" equals "skuba node upgrade apply --user sles --target " plus Master Node IPS
    And I run UPGRADE VARS:"upgradeapply" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "to date"
    And I run UPGRADE VARS:"upgradeapply" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "to date"
    And I run UPGRADE VARS:"upgradeapply" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "to date"


    