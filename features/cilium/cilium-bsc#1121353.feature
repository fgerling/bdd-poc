#TC description: https://github.com/SUSE/caasp-test-cases/pull/22/files#diff-6b705e1a52bb5b33c1c2efee4a329c2f
#BUG: https://bugzilla.suse.com/show_bug.cgi?id=1121353
#PR: https://github.com/SUSE/skuba/pull/967

Feature: bsc#1121353 - Kubernetes – Master node pod configured with Privileged PSP

Scenario: Checking if Privileged Pods
    Given "skuba" exist in gopath
    And VARIABLE "work-folder" I get from CONFIG
    When I run "skuba cluster status" in VAR:"work-folder" directory
    Then the output contains "master" and "worker"
    When I run "kubectl get all --namespace=kube-system"
    Then the output contains "cilium" and "dex"

    When I run "kubectl get pods --namespace=kube-system"
    When VARIABLE "privileged-pods" equals ContainersFROMOutput "cilium-"
    And VARIABLES "commandchecks" equals "kubectl describe pod -n kube-system " plus VAR:"privileged-pods"
    And I run VARS:"commandchecks" and check for "psp" and "kubernetes.io/psp: suse.caasp.psp.privileged"

    When I run "kubectl get pods --namespace=kube-system"
    When VARIABLE "privileged-pods" equals ContainersFROMOutput "kube-proxy-"
    And VARIABLES "commandchecks" equals "kubectl describe pod -n kube-system " plus VAR:"privileged-pods"
    And I run VARS:"commandchecks" and check for "psp" and "kubernetes.io/psp: suse.caasp.psp.privileged"

    When I run "kubectl get pods --namespace=kube-system"
    When VARIABLE "privileged-pods" equals ContainersFROMOutput "kured-"
    And VARIABLES "commandchecks" equals "kubectl describe pod -n kube-system " plus VAR:"privileged-pods"
    And I run VARS:"commandchecks" and check for "psp" and "kubernetes.io/psp: suse.caasp.psp.privileged"
