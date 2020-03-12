package features

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (test *TestRun) IRunSkubaUpgradePlanFirstMasterInVARDirectory(arg1 string) error {
	var firstmaster string
	for key, _ := range test.VarMap {
		if strings.Contains(key, "master") && strings.Contains(key, "00") {
			firstmaster = key
		}
	}
	cmd := []string{"skuba", "upgrade", "plan", firstmaster}
	output, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		test.Err = err
		fmt.Fprintf(os.Stdout, "error: %v", err)
	}
	test.Output = output
	return nil
}
