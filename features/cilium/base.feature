Feature: cilium-basic
  Scenario: cilium is properly deployed and working
    Given In namespace "kube-system" DaemonSet "cilium" exists
    Then In namespace "kube-system" DaemonSet "cilium" should be ready
    Given In namespace "kube-system" Deployment "cilium-operator" exists
    Then In namespace "kube-system" Deployment "cilium-operator" should be ready

    When I exec with kubectl in namespace "kube-system" in pod "ds/cilium" the command "cilium version"
    Then the output contains "Client: 1.6.6"
    And the output contains "Daemon: 1.6.6"

  Scenario: no leftovers from a migration
    Given In namespace "kube-system" no pods with labels "k8s-app=cilium-pre-flight-check" exist

  Scenario: cilium uses CRD instead of etcd
    Given In namespace "kube-system" ConfigMap "cilium-config" exists
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
