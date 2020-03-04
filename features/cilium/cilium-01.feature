# doc1:   https://gitlab.suse.de/mkravec/scripts/blob/master/tests/cilium.sh
# doc2:   http://docs.cilium.io/en/v1.6/gettingstarted/http/
# TC:   https://github.com/fgerling/bdd-poc
# This is a basic test for cilium (no PR or BSC provided)

Feature: cilium-basic

 Scenario: Test-Cilium-Basic on Skuba Cluster
    Given "skuba" exist in gopath
    #And VARIABLE "work-folder" equals "/Users/alexeitighineanu/go/src/github.com/fgerling/bdd-poc/imba-cluster"
    When I run "skuba cluster status"       
    Then the output contains "master" and "worker"
    When I run "kubectl get all --namespace=kube-system"
    Then the output contains "cilium" and "dex"

    Scenario: Deploy the starwars cilium pods
    When I run "kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/http-sw-app.yaml"
    And wait "10 seconds"
    When I run "kubectl get pods --selector=org=empire" 
    Then the output contains "deathstar" and "tiefighter"
    When I run "kubectl get pods --selector=org=alliance"

    Scenario: Check the starwars cilium pods
    When I run "kubectl get pods"                  
    And grep for "xwing"
    Then the output contains "running" and "running"
    And grep for "deathstar"
    Then the output contains "running" and "running"
    And grep for "tiefighter"
    Then the output contains "running" and "running"

    Scenario: Test number1 if empire's ship is allowed into empire space
    And I run "kubectl exec tiefighter -- curl -sm10 -XPOST deathstar.default.svc.cluster.local/v1/request-landing"
    Then the output contains "Ship" and "landed"

    When I run "kubectl exec xwing -- curl -sm10 -XPOST deathstar.default.svc.cluster.local/v1/request-landing"
    Then the output contains "Ship" and "landed"

    Scenario: Test number2 if policies work properly
    When I run "kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/sw_l3_l4_policy.yaml"
    And I run "kubectl exec tiefighter -- curl -sm10 -XPOST deathstar.default.svc.cluster.local/v1/request-landing"
    Then the output contains "Ship" and "landed"

    When I run "kubectl exec xwing -- curl -sm10 -XPOST deathstar.default.svc.cluster.local/v1/request-landing" expecting ERROR
    And wait "10 seconds"
    Then the error contains "exit" and "28"

    Scenario: Inspecting the policies
    When I run "kubectl -n kube-system get pods -l k8s-app=cilium"
    Then the output contains "cilium-" and "running"
    When VARIABLE "cilium-container" equals ContainerFROMOutput "cilium"
    And VARIABLE "command5" equals "kubectl -n kube-system exec " plus VAR:"cilium-container" plus " -- cilium endpoint list"
    And I run VAR:"command5"
    And grep for "class=deathstar"
    Then the output contains "enabled" and "disabled"
    And I run VAR:"command5"
    And grep for "class=xwing"
    Then the output contains "disabled" and "disabled"
    And I run VAR:"command5"
    And grep for "class=tiefighter"
    Then the output contains "disabled" and "disabled"
    And I run "kubectl get cnp"
    When VARIABLE "cilium-rule" equals ContainerFROMOutput "rule"
    And VARIABLE "command6" equals "kubectl describe cnp " plus VAR:"cilium-rule"
    And grep for "Class:" 
    Then the output contains "deathstar" and "deathstar"
    And grep for "Org:"
    Then the output contains "empire" and "empire"
    And grep for "Description:"
    Then the output contains "policy to restrict" and "empire ships only"

    Scenario: Applying new policy for exhaust port
    When I run "kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/sw_l3_l4_l7_policy.yaml"
    Then the output contains "ciliumnetworkpolicy" and "configured"
    And wait "1 seconds"
    #When I run "kubectl exec tiefighter -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing"
    #Then the output contains "Ship" and "landed"
    #When I run "kubectl exec tiefighter -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/exhaust-port"
    #Then the output contains "Access" and "denied"


    Scenario: Deleting the policies
    When I run "kubectl delete -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/sw_l3_l4_l7_policy.yaml"
    Then the output contains "ciliumnetworkpolicy" and "deleted"

    Scenario: Deleting the pods
    When VARIABLE "work-folder" equals "/Users/alexeitighineanu/go/src/github.com/fgerling/bdd-poc/imba-cluster"
    When I run "kubectl delete -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/http-sw-app.yaml"
    Then the output contains "deathstar" and "deleted"
    Then the output contains "xwing" and "deleted"
    Then the output contains "tiefighter" and "deleted"