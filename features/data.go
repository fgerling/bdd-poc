package features

type TestRun struct {
	Output []byte
	VarMap map[string]string
	Err    error
}
