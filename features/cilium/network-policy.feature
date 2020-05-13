Feature: cilium-network-policies
  Scenario: Testing pods are ready
    Given httpbin must be ready
    And tblshoot must be ready

  Scenario: There are no network policies
#    Given I delete all "CiliumNetworkPolicy" resources in "default" namespace
#    Then there is no "CiliumNetworkPolicy" resource in "default" namespace
    Given I run "kubectl delete cnp --all"

  Scenario: http traffic is allowed without network policy
    When I send "<method>" request to "<path>"
    Then the output contains "<retcode>"

    Examples:
      | method | path                  | retcode |
      | DELETE | /anything/allowed     | 200     |
      | DELETE | /anything/not-allowed | 200     |
      | GET    | /anything/allowed     | 200     |
      | GET    | /anything/not-allowed | 200     |
      | PUT    | /anything/allowed     | 200     |
      | PUT    | /anything/not-allowed | 200     |

  Scenario: dns traffic is allowed without network policy
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


  Scenario: Apply a default network policy to deny all traffic
    Given I run "kubectl apply -f manifests/cilium/network-policy-deny-all.yaml"

  Scenario: dns traffic is forbidden by the default network policy
    When I resolve "<name>" and fails
    Then the output contains "connection timed out"
    When I reverse resolve "<ip>" and fails
    Then the output contains "connection timed out"

    Examples:
      | name            | ip        |
      | one.one.one.one | 1.1.1.1   |
      | one.one.one.one | 1.0.0.1   |
      | dns.google      | 8.8.8.8   |
      | localhost       | 127.0.0.1 |


  Scenario: Apply network policy to allow DNS traffic
    Given I run "kubectl apply -f manifests/cilium/network-policy-allow-dns.yaml"

  Scenario: dns is allowed by network policy
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


  Scenario: http traffic is still not allowed by the default network policy
    When I send "<method>" request to "<path>" and fails
    Then the output contains "<retcode>"

    Examples:
      | method | path                  | retcode |
      | DELETE | /anything/allowed     | 000     |
      | GET    | /anything/allowed     | 000     |
      | PUT    | /anything/allowed     | 000     |


  Scenario: Apply network policy to allow http traffic at layer 3
    Given I run "kubectl apply -f manifests/cilium/network-policy-allow-http-l3.yaml"

  Scenario: http traffic is allowed by network policy at layer 3
    When I send "<method>" request to "<path>"
    Then the output contains "<retcode>"

    Examples:
      | method | path                  | retcode |
      | DELETE | /anything/allowed     | 200     |
      | DELETE | /anything/not-allowed | 200     |
      | GET    | /anything/allowed     | 200     |
      | GET    | /anything/not-allowed | 200     |
      | PUT    | /anything/allowed     | 200     |
      | PUT    | /anything/not-allowed | 200     |


  Scenario: Apply network policy to allow http traffic with Layer 7 filtering
    Given I run "kubectl delete -f manifests/cilium/network-policy-allow-http-l3.yaml"
    Then I run "kubectl apply -f manifests/cilium/network-policy-allow-delete-get-put-l7.yaml"

  Scenario: http traffic is filtered by network policy at layer 7
    When I send "<method>" request to "<path>"
    Then the output contains "<retcode>"

    Examples:
      | method | path                  | retcode |
      | DELETE | /anything/allowed     | 200     |
      | DELETE | /anything/not-allowed | 403     |
      | GET    | /anything/allowed     | 200     |
      | GET    | /anything/not-allowed | 403     |
      | PUT    | /anything/allowed     | 200     |
      | PUT    | /anything/not-allowed | 403     |
