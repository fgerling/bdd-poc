# TC:   https://github.com/SUSE/caasp-test-cases
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

    When I run "skuba addon upgrade apply" in VAR:"imba-cluster" directory
    Then the output contains "congratulations" or "not all"

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

    When VARIABLES "upgradeapply" equals "skuba node upgrade apply --user sles --sudo --target " plus Master Node IPS
    And I run UPGRADE VARS:"upgradeapply" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "to date" or "there are addon upgrades available"
    And wait "30 seconds"
    And I run UPGRADE VARS:"upgradeapply" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "to date" or "there are addon upgrades available"
    And wait "30 seconds"
    And I run UPGRADE VARS:"upgradeapply" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "to date" or "there are addon upgrades available"
    And wait "30 seconds"

# UPGRADING THEN WORKERS 
    #Scenario: Upgrading Workers
    When VARIABLES "commandupgrades2" equals "skuba node upgrade plan " plus Worker Nodes
    #And VARIABLE "imba-cluster" equals "/root/go/src/github.com/fgerling/bdd-poc/imba-cluster"
    And I run UPGRADE VARS:"commandupgrades2" in VAR:"imba-cluster" directory
    Then the output contains "kubelet" and "cri-o"
    And I run UPGRADE VARS:"commandupgrades2" in VAR:"imba-cluster" directory
    Then the output contains "kubelet" and "cri-o"

    When VARIABLES "upgradeapply2" equals "skuba node upgrade apply --user sles --sudo --target " plus Worker Node IPS
    And I run UPGRADE VARS:"upgradeapply2" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "to date" or "there are addon upgrades available"
 
    When I run "skuba addon upgrade apply" in VAR:"imba-cluster" directory
    Then the output contains "not all nodes" or "successfully" or "congratulations"

    When VARIABLES "upgradeapply3" equals "skuba node upgrade apply --user sles --sudo --target " plus Worker Node IPS
    And I run UPGRADE VARS:"upgradeapply3" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "to date" or "there are addon upgrades available"

    When I run "skuba addon upgrade apply" in VAR:"imba-cluster" directory
    Then the output contains "successfully" or "congratulations"
