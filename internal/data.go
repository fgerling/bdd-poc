package features

type TestRun struct {
	Output       []byte
	VarMap       map[string]string
	UpgradeCheck map[string]NodeCheck
	Err          error
}

type NodeCheck struct {
	PlanDone bool
	UPDone   bool
	IP       string
}
