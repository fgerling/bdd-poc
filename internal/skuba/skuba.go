package skuba

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/cucumber/godog"
)

func SkubaHasVersion(version string) error {
	cmd := exec.Command("skuba", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	if !strings.Contains(string(output), version) {
		return errors.New("Version does not match")
	}
	return nil
}

func ClusterStatus(key, value string) error {
	cmd := exec.Command("skuba", "cluster", "status")
	cmd.Dir = "cluster"
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return err
}
func IUpgradeTheClusterWithSkubaReleaseVersion(arg1 string) error {
	return godog.ErrPending
}
