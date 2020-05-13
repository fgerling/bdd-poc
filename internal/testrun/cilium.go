package testrun

import (
	"github.com/cucumber/messages-go/v10"
)

var (
	ciliumConfigMap = map[string]string{
		"Name":      "cilium-config",
		"Namespace": "kube-system",
	}
)

// CiliumConfigMapDoesHaveTheOptions returns if cilium ConfigMap does have the key/value
// has the key/value provided in the content of a PickleStepArgument_PickleDocString
func (t *TestRun) CiliumConfigMapDoesHaveTheOptions(options *messages.PickleStepArgument_PickleDocString) error {
	return t.ConfigMapDoesHaveTheOptions(ciliumConfigMap["Namespace"], ciliumConfigMap["Name"], options.Content)
}

// CiliumConfigMapDoesNotHaveTheOptions returns if cilium ConfigMap
// has the key/value provided in the content of a PickleStepArgument_PickleDocString
func (t *TestRun) CiliumConfigMapDoesNotHaveTheOptions(options *messages.PickleStepArgument_PickleDocString) error {
	return t.ConfigMapDoesNotHaveTheOptions(ciliumConfigMap["Namespace"], ciliumConfigMap["Name"], options.Content)
}
