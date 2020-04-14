Feature: cilium-basic
  Scenario: cilium is properly deployed and working
    Given DaemonSet "cilium" in namespace "kube-system" exists
    Then DaemonSet "cilium" in namespace "kube-system" should be ready
    Given Deployment "cilium-operator" in namespace "kube-system" exists
    Then Deployment "cilium-operator" in namespace "kube-system" should be ready

  Scenario: cilium uses CRD instead of etcd
    Given ConfigMap "cilium-config" in namespace "kube-system" exists
#    Then cilium ConfigMap does have the keys and values:
#      | key                      | value |
#      | identity-allocation-mode | crd   |
#      | preallocate-bpf-maps     | false |
    Then cilium ConfigMap does have the options:
      """
      {
        "bpf-ct-global-any-max": "262144",
        "bpf-ct-global-tcp-max": "524288",
        "debug": "false",
        "enable-ipv4": "true",
        "enable-ipv6": "false",
        "identity-allocation-mode": "crd",
        "preallocate-bpf-maps": "false"
      }
      """
    Then cilium ConfigMap does not have the options:
      """
      {
        "etcd-config": "",
        "kvstore": "",
        "kvstore-opt": ""
      }
      """

    #And ConfigMap "cilium-config" in namespace "kube-system" does not have the keys:

#    When I do test "cat /etc/hostname"
#    Then the output contains "010084073204"
#    Then the output shoud match the output the command "kubectl get nodes -o go-template --template='{{len .items}}'"
#    When I run "kubectl -n kube-system wait pods --for=condition=Ready --timeout=120s -l k8s-app=kube-dns -ojsonpath='{.status.containerStatuses[0].ready}'"
#    Then the output contains "true"
#    When I run "kubectl -n kube-system get deploy coredns -ojsonpath={.status.readyReplicas}"
#    Then the output contains "2"
