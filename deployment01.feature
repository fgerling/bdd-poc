# doc: https://github.com/SUSE/skuba/blob/master/README.md
# TO FIX: add a ${bdd-poc_rootDir}
#         add a ${kube_admin.conf rootDir}


Feature: kubernetes deployments

  Scenario: Create deployment on CaaSPv4
    Given there is "imba-cluster" directory
    And "skuba" exist in gopath
    When I run "kubectl get all --namespace=kube-system" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "cilium" and "dex"
    When I run "kubectl create -f deployment01.yaml" in "/home/atighineanu/golang/src/bdd-poc" directory
    Then the output contains "deployment01" and "created"
    And wait "30 seconds"
    When I run "kubectl get deployments" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "deployment01" and "3/3"
  
  Scenario: Scale the deployment on CaaSPv4
    Given there is "imba-cluster" directory
    When I run "kubectl get all --namespace=kube-system" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "cilium" and "dex"
    And wait "5 seconds"
    When I run "kubectl get deployments" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "deployment01" and "3/3"
    When I run "kubectl scale deployment deployment01 --replicas=10" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "deployment01" and "scaled"
    And wait "30 seconds"
    When I run "kubectl describe deployments deployment01" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    And grep for "replicas:"
    Then the output contains "10 desired" and "10 total"

 Scenario: Label, change, edit deployment on CaaSPv4
    Given there is "imba-cluster" directory
    When I run "kubectl get all --namespace=kube-system" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "cilium" and "dex"
    And wait "5 seconds"
    When I run "kubectl get deployments" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "deployment01" and "10/10"
    And wait "1 seconds"
    When I run "kubectl label deployments deployment01 environment=premium" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "deployment01" and "labeled"
    When I run "kubectl describe deployments deployment01" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    And grep for "environment="
    Then the output contains "environment" and "premium"