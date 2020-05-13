package testrun

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type TestRun struct {
	CombinedOutput []byte
	StdOut         []byte
	StdErr         []byte
	VarMap         map[string]string
	UpgradeCheck   map[string]NodeCheck
	Err            error
	Config         Config

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
