package features

type TestRun struct {
	Output       []byte
	VarMap       map[string]string
	UpgradeCheck map[string]bool
	Err          error
}
