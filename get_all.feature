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





# OUTPUT EXAMPLE:
#   Scenario: get all resources                                                                                             # get_all.feature:5
#    Given there is "imba-cluster" directory                                                                               # makefile_test.go:52 -> theDirectoryExsist
#    And "skuba" exist in gopath                                                                                           # makefile_test.go:14 -> existInGopath
#    When I run "kubectl get all --namespace=kube-system" in "/home/atighineanu/golang/src/bdd-poc/imba-cluster" directory # makefile_test.go:67 -> iRunInDirectory
#    Then the output contains "cilium" and "dex"                                                                           # makefile_test.go:80 -> theOutputContainsAnd
#    And the output contains "replicaset.apps/cilium-operator" and "replicaset.apps/coredns"                               # makefile_test.go:80 -> theOutputContainsAnd
#    And the output contains "replicaset.apps/oidc-dex" and "replicaset.apps/oidc-gangway"                                 # makefile_test.go:80 -> theOutputContainsAnd
#    And the output contains "daemonset.apps/cilium" and "daemonset.apps/kube-proxy"                                       # makefile_test.go:80 -> theOutputContainsAnd
#    And the output contains "daemonset.apps/kured" and "deployment.apps/cilium-operator" 
#1 scenarios (1 passed)
#8 steps (8 passed)
#2.912988005s