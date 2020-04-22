Feature: coredns-basic
  Background:
    Given the environment variable "KUBECONFIG" is set to "cluster/admin.conf"

  Scenario: the 2 coredns pods are ready
    When I run "kubectl -n kube-system get deploy coredns -ojsonpath={.spec.replicas}"
    Then the output contains "2"
    When I run "kubectl -n kube-system wait pods --for=condition=Ready --timeout=120s -l k8s-app=kube-dns -ojsonpath='{.status.containerStatuses[0].ready}'"
    Then the output contains "true"
    When I run "kubectl -n kube-system get deploy coredns -ojsonpath={.status.readyReplicas}"
    Then the output contains "2"

  Scenario: kube-dns service exists
    When I run "kubectl -n kube-system get svc kube-dns -ojsonpath={.metadata.name}"
    Then the output contains "kube-dns"

  Scenario: testing pod dnsutils must be ready
    When I run "kubectl apply -f https://raw.githubusercontent.com/lcavajani/cagibi/master/kubernetes/dnsutils/dnutils-pod.yaml"
    Then the output contains "pod/dnsutils"
    And the output contains "created" or "unchanged"
    When I run "kubectl wait --for=condition=Ready --timeout=120s pod/dnsutils"
    Then the output contains "pod/dnsutils condition met"

  Scenario: can resolve kubernetes api IP from its internal FQDN
    When I run "kubectl exec dnsutils dig +short kubernetes.default.svc.cluster.local"
    Then the output shoud match the output the command "kubectl get svc kubernetes -ojsonpath={.spec.clusterIP}"

  Scenario: can reverse resolve kubernetes api FQDN from its internal IP
    When I run "kubectl exec dnsutils -- dig +short -x 10.96.0.1"
    Then the output contains "kubernetes.default.svc.cluster.local."

  Scenario: can resolve IP from external FQDN
    When I run "kubectl exec dnsutils dig +short one.one.one.one"
    Then the output contains "1.1.1.1"
    When I run "kubectl exec dnsutils dig +short 127.0.0.1.omg.howdoi.website"
    Then the output contains "127.0.0.1"
    When I run "kubectl exec dnsutils dig +short 8.8.8.8.omg.howdoi.website"
    Then the output contains "8.8.8.8"

  Scenario: can reverse resolve FQDN from external IP
    When I run "kubectl exec dnsutils -- dig +short -x 1.1.1.1"
    Then the output contains "one.one.one.one."
    When I run "kubectl exec dnsutils -- dig +short -x 8.8.8.8"
    Then the output contains "dns.google."

  Scenario: testing pod dnsutils-netcat must be ready
    When I run "kubectl apply -f https://raw.githubusercontent.com/lcavajani/cagibi/master/kubernetes/dnsutils/dnutils-netcat-pod.yaml"
    Then the output contains "pod/dnsutils-netcat"
    And the output contains "created" or "unchanged"
    When I run "kubectl wait --for=condition=Ready --timeout=120s pod/dnsutils-netcat"
    Then the output contains "pod/dnsutils-netcat condition met"

  Scenario: can create dns A and PTR entry for service
    When I run "kubectl apply -f https://raw.githubusercontent.com/lcavajani/cagibi/master/kubernetes/dnsutils/dnutils-netcat-service.yaml"
    Then the output contains "service/dnsutils-netcat"
    And the output contains "created" or "unchanged"
    When I run "kubectl get svc dnsutils-netcat -ojsonpath={.spec.clusterIP}"
    Then the output contains a valid ip address
    When I run "kubectl exec dnsutils -- dig +short dnsutils-netcat.default.svc.cluster.local"
    Then the output shoud match the output the command "kubectl get po dnsutils -ojsonpath={.spec.clusterIP}"
    When I run "kubectl exec dnsutils -- dig +short -x 10.100.100.100"
    Then the output contains "dnsutils-netcat.default.svc.cluster.local"

  Scenario: can create dns A and PTR entry for headless service
    When I run "kubectl apply -f https://raw.githubusercontent.com/lcavajani/cagibi/master/kubernetes/dnsutils/dnutils-netcat-service-headless.yaml"
    Then the output contains "service/dnsutils-netcat-headless"
    And the output contains "unchanged" or "created"
    When I run "kubectl get svc dnsutils-netcat-headless -ojsonpath={.spec.clusterIP}"
    Then the output contains "None"
    When I run "kubectl exec dnsutils -- dig +short dnsutils-netcat-headless.default.svc.cluster.local"
    Then the output shoud match the output the command "kubectl get po dnsutils-netcat -ojsonpath={.spec.clusterIP}"
