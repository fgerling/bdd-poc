package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/cucumber/messages-go/v10"
)

var (
	kindAppsReplicaSet = schema.GroupKind{Group: "apps", Kind: "ReplicaSet"}
	kindAppsDeployment = schema.GroupKind{Group: "apps", Kind: "Deployment"}
	kindAppsDaemonSet  = schema.GroupKind{Group: "apps", Kind: "DaemonSet"}
	kindBatchJob       = schema.GroupKind{Group: "batch", Kind: "Job"}
	kindConfigMap      = schema.GroupKind{Group: "core", Kind: "ConfigMap"}
	kindSecret         = schema.GroupKind{Group: "core", Kind: "Secret"}
	kindService        = schema.GroupKind{Group: "core", Kind: "Service"}

	ciliumConfigMap = map[string]string{
		"Name":      "cilium-config",
		"Namespace": "kube-system",
	}
)

type TestRun struct {
	Output       []byte
	VarMap       map[string]string
	UpgradeCheck map[string]NodeCheck
	Err          error
	Config       Config
	// Kubernetes related
	RestConfig *rest.Config
	ClientSet  *kubernetes.Clientset
	//	Resource   interface {
	//		IsReady() error
	//	}
}

type NodeCheck struct {
	PlanDone bool
	UPDone   bool
	IP       string
}

type Config struct {
	ClusterDir string `json:"ClusterDir"`
}

func GetKubeconfigPath() (string, error) {
	// Read KUBECONFIG only from env var
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		return "", fmt.Errorf("KUBECONFIG variable not exported")
	}

	return kubeconfig, nil
}

func GetClientRestConfig() (*rest.Config, error) {
	kubeconfig, err := GetKubeconfigPath()
	if err != nil {
		log.Fatal(err)
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// create the clientSet
func (t *TestRun) CreateClientSet() (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(t.RestConfig)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func (t *TestRun) GetPod(name, ns string) (*v1.Pod, error) {
	pod, err := t.ClientSet.CoreV1().Pods(ns).Get(name, metav1.GetOptions{})
	//if !(errors.IsNotFound(err)) {
	//	fmt.Printf("Found pod %s in namespace %s\n", name, ns)
	//}
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (t *TestRun) GetDeployment(name, ns string) (*appsv1.Deployment, error) {
	dp, err := t.GetRuntimeObjectForKind(kindAppsDeployment, name, ns)
	if err != nil {
		return nil, err
	}

	return dp.(*appsv1.Deployment), nil
}

func (t *TestRun) GetDaemonset(name, ns string) (*appsv1.DaemonSet, error) {
	ds, err := t.GetRuntimeObjectForKind(kindAppsDaemonSet, name, ns)
	if err != nil {
		return nil, err
	}

	return ds.(*appsv1.DaemonSet), nil
}

// DaemonSetInNamespaceExists returns an error if DaemonSet does not exist
func (t *TestRun) DaemonSetInNamespaceExists(name, ns string) error {
	_, err := t.GetDaemonset(name, ns)
	return err
}

// DaemonSetInNamespaceIsReady returns an error if DaemonSet is not ready
func (t *TestRun) DaemonSetInNamespaceIsReady(name, ns string) error {
	ds, err := t.GetRuntimeObjectForKind(kindAppsDaemonSet, name, ns)
	if err != nil {
		return err
	}

	return IsRuntimeObjectReady(ds)
}

// DeploymentInNamespaceExists returns an error if Deployment does not exist
func (t *TestRun) DeploymentInNamespaceExists(name, ns string) error {
	//	_, err := t.GetRuntimeObjectForKind(kindAppsDeployment, name, ns)
	_, err := t.GetDeployment(name, ns)
	return err
}

// DeploymentInNamespaceIsReady returns an error if Deployment is not ready
func (t *TestRun) DeploymentInNamespaceIsReady(name, ns string) error {
	dp, err := t.GetRuntimeObjectForKind(kindAppsDeployment, name, ns)
	if err != nil {
		return err
	}

	return IsRuntimeObjectReady(dp)
}

func (t *TestRun) GetConfigMap(name, ns string) (*v1.ConfigMap, error) {
	cm, err := t.GetRuntimeObjectForKind(kindConfigMap, name, ns)
	if err != nil {
		return nil, err
	}

	return cm.(*v1.ConfigMap), nil
}

func (t *TestRun) ConfigMapInNamespaceExists(name, ns string) error {
	_, err := t.GetConfigMap(name, ns)
	return err
}

func (t *TestRun) ciliumConfigMapDoesHaveTheOptions(options *messages.PickleStepArgument_PickleDocString) error {
	var expected map[string]string
	c, _ := t.GetConfigMap(ciliumConfigMap["Name"], ciliumConfigMap["Namespace"])

	if err := json.Unmarshal([]byte(options.Content), &expected); err != nil {
		return err
	}

	for k, v := range expected {
		if c.Data[k] != v {
			return fmt.Errorf("incorrect option in current config, %v: %v", k, v)
		}
	}

	return nil
}

func (t *TestRun) ciliumConfigMapDoesNotHaveTheOptions(options *messages.PickleStepArgument_PickleDocString) error {
	var expected map[string]string
	c, err := t.GetConfigMap(ciliumConfigMap["Name"], ciliumConfigMap["Namespace"])
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(options.Content), &expected); err != nil {
		return err
	}

	for k, _ := range expected {
		if _, err := c.Data[k]; err {
			return fmt.Errorf("non expected options exists in current config, %v", k)
		}
	}

	return nil
}

// GetRuntimeObjectForKind returns a runtime.Object based on its GroupKind,
// name and namespace
func (t *TestRun) GetRuntimeObjectForKind(kind schema.GroupKind, name, ns string) (runtime.Object, error) {
	switch kind {
	case kindAppsReplicaSet:
		return t.ClientSet.AppsV1().ReplicaSets(ns).Get(name, metav1.GetOptions{})
	case kindAppsDeployment:
		return t.ClientSet.AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	case kindAppsDaemonSet:
		return t.ClientSet.AppsV1().DaemonSets(ns).Get(name, metav1.GetOptions{})
	case kindBatchJob:
		return t.ClientSet.BatchV1().Jobs(ns).Get(name, metav1.GetOptions{})
	case kindConfigMap:
		return t.ClientSet.CoreV1().ConfigMaps(ns).Get(name, metav1.GetOptions{})
	default:
		return nil, fmt.Errorf("Unsupported kind when getting runtime object: %v", kind)
	}
}

// IsRuntimeObjectReady returns if a given object is ready
func IsRuntimeObjectReady(obj runtime.Object) error {
	switch typed := obj.(type) {
	case *appsv1.ReplicaSet:
		if typed.Status.Replicas >= typed.Status.ReadyReplicas {
			return nil
		}
		return fmt.Errorf("Some pods are not ready %d/%d", typed.Status.Replicas, typed.Status.ReadyReplicas)
	case *appsv1.Deployment:
		if typed.Status.Replicas >= typed.Status.ReadyReplicas {
			return nil
		}
		return fmt.Errorf("Some pods are not ready %d/%d", typed.Status.Replicas, typed.Status.ReadyReplicas)
	case *appsv1.DaemonSet:
		if typed.Status.DesiredNumberScheduled >= typed.Status.NumberReady {
			return nil
		}
		return fmt.Errorf("Some pods are not ready %d/%d", typed.Status.DesiredNumberScheduled, typed.Status.NumberReady)
	default:
		return fmt.Errorf("Unsupported kind when getting number of replicas: %v", obj)
	}
}

func (t *TestRun) GetNode(nodeName, namespace string) (*v1.Node, error) {
	node, err := t.ClientSet.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (t *TestRun) ExecuteCommandInPod(namespace, podName, containerName string, cmd []string) ([]byte, []byte, error) {
	const tty = false
	var stdin io.Reader
	var stdout, stderr bytes.Buffer

	req := t.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName)
		//		Param("container", pod.Spec.Containers[0].Name)

	options := &v1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       tty,
	}

	req.VersionedParams(options, scheme.ParameterCodec)
	err := t.executeCmd("POST", req.URL(), stdin, &stdout, &stderr, tty)

	//if options.PreserveWhitespace {
	//	return stdout.String(), stderr.String(), err
	//}
	//return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
	return stdout.Bytes(), stderr.Bytes(), err
}

func (t *TestRun) executeCmd(method string, url *url.URL, stdin io.Reader, stdout, stderr io.Writer, tty bool) error {
	exec, err := remotecommand.NewSPDYExecutor(t.RestConfig, "POST", url)
	if err != nil {
		return err
	}

	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    tty,
	})
}

//type Deployment struct {
//	*appsv1.Deployment
//}
//
//func (d Deployment) IsReady() error {
//	if d.Status.Replicas == d.Status.ReadyReplicas {
//		return nil
//	}
//
//	return errors.New(fmt.Sprintf("Some pods are not ready %d/%d", d.Status.Replicas, d.Status.ReadyReplicas))
//}

func IDoTest(arg1 string) error {
	var err error
	var cmd []string
	//pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	var t TestRun
	t.RestConfig, err = GetClientRestConfig()
	if err != nil {
		log.Fatal(err)
	}

	t.ClientSet, err = t.CreateClientSet()
	if err != nil {
		log.Fatal(err)
	}

	namespace := "kube-system"
	podName := "cilium-cdptx"
	//	cmd = append(cmd, "cat", "/etc/hostname")
	//cmd = append(cmd, "echo", "toto")
	cmd = strings.Split(arg1, " ")
	// get pod
	pod, err := t.GetPod(podName, namespace)
	if err != nil {
		log.Fatal(err)
	}

	stdout, _, err := t.ExecuteCommandInPod(namespace, podName, pod.Spec.Containers[0].Name, cmd)
	if err != nil {
		log.Fatal(err)
	}

	Output = stdout

	return nil
}
