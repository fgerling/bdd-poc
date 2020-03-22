package features

type TestRun struct {
	Output       []byte
	VarMap       map[string]string
	UpgradeCheck map[string]NodeCheck
	Err          error
	Config       Config
}

type NodeCheck struct {
	PlanDone bool
	UPDone   bool
	IP       string
}

type Config struct {
	ClusterDir string `json:"ClusterDir"`
}
