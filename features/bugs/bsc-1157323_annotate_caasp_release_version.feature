# BUG: https://bugzilla.suse.com/show_bug.cgi?id=1157323
# PR: https://github.com/SUSE/skuba/pull/894
# DOC PR: https://github.com/SUSE/doc-caasp/pull/635

Feature: bsc#1157323 Annotate CaaSP release version
  Scenario Template: Deploy a cluster
    Given skuba has version "<skuba-version>"
    And I have a cluster "cluster"
    When I execute "skuba cluster status" in "cluster"
    Then skuba cluster status "CAASP-RELEASE-VERSION" is "<caasp-release-version>"

    Examples:
      | skuba-version | caasp-release-version |
      | 1.1.2         | 4.0.3                 |
      | 1.2.1         | 4.1.0                 |
      | untagged      | <none>                |

  Scenario Template: Upgrade a cluster
    Given skuba has version "1.1.1"
    And I have a cluster "cluster"
    When I upgrade the cluster with skuba release version "<new-skuba-version>"
    And I execute "skuba cluster status"
    Then skuba cluster status "CAASP-RELEASE-VERSION" is "<caasp-release-version>"

    Examples:
      | new-skuba-version | caasp-release-version |
      | 1.1.2             | 4.0.3                 |
      | 1.2.1             | 4.1.0                 |
