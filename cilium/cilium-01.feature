# doc1:   https://gitlab.suse.de/mkravec/scripts/blob/master/tests/cilium.sh
# doc2:   http://docs.cilium.io/en/v1.6/gettingstarted/http/
# doc3:   https://github.com/fgerling/bdd-poc

Feature: cilium-basic

 Scenario: Test-Cilium-Basic on Skuba Cluster
    Given there is "imba-cluster" directory
    And "skuba" exist in gopath
    And VARIABLE "work-folder" equals "/home/atighineanu/golang/src/bdd-poc/imba-cluster"
    When I run "skuba cluster status" in VAR:"work-folder" directory        
    Then the output contains "master" and "worker"
    When I run "kubectl get all --namespace=kube-system" in VAR:"work-folder" directory
    Then the output contains "cilium" and "dex"
    And VARIABLE "curlreq" equals "curl -sm10 -XPOST deathstar.default.svc.cluster.local/v1/request-landing" 

    Scenario: Deploy the starwars cilium pods
    When I run "kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/http-sw-app.yaml" in VAR:"work-folder" directory
    And wait "10 seconds"
    When I run "kubectl get pods --selector=org=empire" in VAR:"work-folder" directory
    Then the output contains "deathstar" and "tiefighter"
    When I run "kubectl get pods --selector=org=alliance" in VAR:"work-folder" directory
    Then the output contains "xwing" and ""

    Scenario: Check the starwars cilium pods
    When I run "kubectl get pods" in VAR:"work-folder" directory                   
    And grep for "xwing"
    Then the output contains "running" and ""
    And grep for "deathstar"
    Then the output contains "running" and ""
    And grep for "tiefighter"
    Then the output contains "running" and ""

    Scenario: Test number1 if empire's ship is allowed into empire space
    When VARIABLE "work-folder" equals "/home/atighineanu/golang/src/bdd-poc/imba-cluster"
    And VARIABLE "curlreq" equals "curl -sm10 -XPOST deathstar.default.svc.cluster.local/v1/request-landing" 
    And VARIABLE "command1" equals "kubectl exec tiefighter -- " plus VAR:"curlreq"
    And I run VAR:"command1" in VAR:"work-folder" directory
    Then the output contains "Ship landed" and ""

    When VARIABLE "command2" equals "kubectl exec xwing -- " plus VAR:"curlreq"
    When I run VAR:"command2" in VAR:"work-folder" directory
    Then the output contains "Ship landed" and ""

    Scenario: Test number2 if policies work properly
    When VARIABLE "work-folder" equals "/home/atighineanu/golang/src/bdd-poc/imba-cluster"
    And VARIABLE "curlreq" equals "curl -sm10 -XPOST deathstar.default.svc.cluster.local/v1/request-landing" 
    When I run "kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/sw_l3_l4_policy.yaml" in VAR:"work-folder" directory 
    And VARIABLE "command3" equals "kubectl exec tiefighter -- " plus VAR:"curlreq"
    And I run VAR:"command3" in VAR:"work-folder" directory
    Then the output contains "Ship landed" and ""

    When VARIABLE "command4" equals "kubectl exec xwing --" plus VAR:"curlreq"
    When I run VAR:"command4" expecting ERROR in VAR:"work-folder" directory
    And wait "10 seconds"
    Then the error contains "exit" and "28"

    Scenario: Inspecting the policies
    When VARIABLE "work-folder" equals "/home/atighineanu/golang/src/bdd-poc/imba-cluster"
    When I run "kubectl -n kube-system get pods -l k8s-app=cilium" in VAR:"work-folder" directory
    Then the output contains "cilium-" and "running"
    When VARIABLE "cilium-container" equals ContainerFROMOutput "cilium"
    And VARIABLE "command5" equals "kubectl -n kube-system exec " plus VAR:"cilium-container" plus " -- cilium endpoint list"
    And I run VAR:"command5" in VAR:"work-folder" directory
    And grep for "class=deathstar"
    Then the output contains "enabled" and "disabled"
    And I run VAR:"command5" in VAR:"work-folder" directory
    And grep for "class=xwing"
    Then the output contains "disabled" and "disabled"
    And I run VAR:"command5" in VAR:"work-folder" directory
    And grep for "class=tiefighter"
    Then the output contains "disabled" and "disabled"
    And I run "kubectl get cnp" in VAR:"work-folder" directory
    When VARIABLE "cilium-rule" equals ContainerFROMOutput "rule"
    And VARIABLE "command6" equals "kubectl describe cnp " plus VAR:"cilium-rule"
    And grep for "Class:" 
    Then the output contains "deathstar" and ""
    And grep for "Org:"
    Then the output contains "empire" and ""
    And grep for "Description:"
    Then the output contains "policy to restrict" and "empire ships only"

    Scenario: Applying new policy for exhaust port
    When VARIABLE "work-folder" equals "/home/atighineanu/golang/src/bdd-poc/imba-cluster"
    When I run "kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/sw_l3_l4_l7_policy.yaml" in VAR:"work-folder" directory
    Then the output contains "ciliumnetworkpolicy" and "configured"
    And wait "1 seconds"
    #When I run "kubectl exec tiefighter -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing" in VAR:"work-folder" directory
    #Then the output contains "Ship" and "landed"
    #When I run "kubectl exec tiefighter -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/exhaust-port" in VAR:"work-folder" directory
    #Then the output contains "Access" and "denied"


    Scenario: Deleting the policies
    When VARIABLE "work-folder" equals "/home/atighineanu/golang/src/bdd-poc/imba-cluster"
    When I run "kubectl delete -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/sw_l3_l4_l7_policy.yaml" in VAR:"work-folder" directory
    Then the output contains "ciliumnetworkpolicy" and "deleted"

    Scenario: Deleting the pods
    When VARIABLE "work-folder" equals "/home/atighineanu/golang/src/bdd-poc/imba-cluster"
    When I run "kubectl delete -f https://raw.githubusercontent.com/cilium/cilium/v1.6/examples/minikube/http-sw-app.yaml" in VAR:"work-folder" directory
    Then the output contains "deathstar" and "deleted"
    Then the output contains "xwing" and "deleted"
    Then the output contains "tiefighter" and "deleted"

   