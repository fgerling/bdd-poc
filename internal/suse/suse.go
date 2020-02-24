package suse

import (
	"os/exec"

	"github.com/cucumber/godog"
)

func MyWorkstationFulfillTheRequirements() error {
	return godog.ErrPending
}

func IInstallThePattern(arg1 string) error {
	return godog.ErrPending
}
func IHaveInPATH(command string) error {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+command)
	return cmd.Run()
}
