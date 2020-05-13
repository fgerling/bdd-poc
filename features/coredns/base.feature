Feature: coredns-basic
  Scenario: Testing pods are ready
    Given tblshoot must be ready

  Scenario: the 2 coredns pods are ready
    Given In namespace "kube-system" Deployment "coredns" exists
    Then In namespace "kube-system" Deployment "coredns" should be ready
    When I run "kubectl -n kube-system get deploy coredns -ojsonpath={.spec.replicas}"
    Then the output contains "2"
    When I run "kubectl -n kube-system get deploy coredns -ojsonpath={.status.readyReplicas}"
    Then the output contains "2"

  Scenario: kube-dns service exists
    When I run "kubectl -n kube-system get svc kube-dns -ojsonpath={.metadata.name}"
    Then the output contains "kube-dns"

  Scenario: can resolve kubernetes api IP from its internal FQDN
    When I resolve "kubernetes.default.svc.cluster.local"
    Then the output shoud match the output the command "kubectl get svc kubernetes -ojsonpath={.spec.clusterIP}"

  Scenario: can reverse resolve kubernetes api FQDN from its internal IP
    When I reverse resolve "10.96.0.1"
    Then the output contains "kubernetes.default.svc.cluster.local"

  Scenario: can resolve external FQDN and reverse resolve external IPs
    When I resolve "<name>"
    Then the output contains "<ip>"
    When I reverse resolve "<ip>"
    Then the output contains "<name>"

    Examples:
      | name            | ip        |
      | one.one.one.one | 1.1.1.1   |
      | one.one.one.one | 1.0.0.1   |
      | dns.google      | 8.8.8.8   |
      | localhost       | 127.0.0.1 |

  Scenario: testing pod dnsutils-netcat must be ready
    When I apply the manifest "manifests/dnsutils/dnsutils-netcat-pod.yaml"
    Then the output contains "pod/dnsutils-netcat"
    And the output contains "created" or "unchanged"
    When I run "kubectl wait --for=condition=Ready --timeout=120s pod/dnsutils-netcat"
    Then the output contains "pod/dnsutils-netcat condition met"

  Scenario: can create dns A and PTR entry for service
    When I apply the manifest "manifests/dnsutils/dnsutils-netcat-service.yaml"
    Then the output contains "service/dnsutils-netcat"
    And the output contains "created" or "unchanged"
    When I run "kubectl get svc dnsutils-netcat -ojsonpath={.spec.clusterIP}"
    Then the output contains a valid ip address
    When I resolve "dnsutils-netcat.default.svc.cluster.local"
    Then the output shoud match the output the command "kubectl get po dnsutils -ojsonpath={.spec.clusterIP}"
    When I reverse resolve "10.100.100.100"
    Then the output contains "dnsutils-netcat.default.svc.cluster.local"

  Scenario: can create dns A and PTR entry for headless service
    When I apply the manifest "manifests/dnsutils/dnsutils-netcat-service-headless.yaml"
    Then the output contains "service/dnsutils-netcat-headless"
    And the output contains "unchanged" or "created"
    When I run "kubectl get svc dnsutils-netcat-headless -ojsonpath={.spec.clusterIP}"
    Then the output contains "None"
    When I resolve "dnsutils-netcat-headless.default.svc.cluster.local"
    Then the output shoud match the output the command "kubectl get po dnsutils-netcat -ojsonpath={.spec.clusterIP}"
