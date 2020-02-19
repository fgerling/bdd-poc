# doc: https://github.com/SUSE/skuba/blob/master/README.md

Feature: kubectl get all 

  Scenario: get all resources
    Given there is "imba-cluster" directory
    And "skuba" exist in gopath
    When I run "kubectl get all --namespace=kube-system" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory
    Then the output contains "cilium" and "dex"
    And the output contains "replicaset.apps/cilium-operator" and "replicaset.apps/coredns"
    And the output contains "replicaset.apps/oidc-dex" and "replicaset.apps/oidc-gangway"
    And the output contains "daemonset.apps/cilium" and "daemonset.apps/kube-proxy"
    And the output contains "daemonset.apps/kured" and "deployment.apps/cilium-operator"