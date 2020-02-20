package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/cucumber/godog"
	suse "github.com/fgerling/bdd-poc/internal/suse"
	git "gopkg.in/src-d/go-git.v4"
)

var Output []byte

func existInGopath(arg1 string) error {
	return theFileExist(path.Join(os.Getenv("GOPATH"), "bin"))
}

func iGitCloneInto(url, target string) error {
	_, err := git.PlainClone(target, false, &git.CloneOptions{
		URL: url,
	})
	return err
}

func iRemoveFromGopath(file string) error {
	return os.Remove(path.Join(os.Getenv("GOPATH"), "bin", file))
}

func iSetTo(variable, value string) error {
	os.Setenv(variable, value)
	return theIsSetTo(variable, value)
}

func theIsSetTo(variable, value string) error {
	if os.Getenv(variable) != value {
		return fmt.Errorf("Env %v is not set to %v", variable, value)
	}
	return nil
}

func theRepositoryExist(repository string) error {
	return theFileExist(path.Join(repository, ".git/"))
}

func theDirectoryExist(dir string) error {
	return theFileExist(dir)
}

func theFileExist(file string) error {
	_, err := os.Stat(file)
	return err
}

func thereIsNoDirectory(target string) error {
	return os.RemoveAll(target)
}

func iRunInDirectory(command, workdir string) error {
	var err error
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workdir
	Output, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return err
}

func theOutputContainsAnd(arg1, arg2 string) error {
	if !strings.Contains(fmt.Sprintf("%s", string(Output)), arg1) && strings.Contains(fmt.Sprintf("%s", string(Output)), arg2) {
		return errors.New("Output does not contain expected arguments")
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^"([^"]*)" exist in gopath$`, existInGopath)
	s.Step(`^I git clone "([^"]*)" into "([^"]*)"$`, iGitCloneInto)
	s.Step(`^I have "([^"]*)" in PATH$`, suse.IHaveInPATH)
	s.Step(`^I install the pattern "([^"]*)"$`, suse.IInstallThePattern)
	s.Step(`^I remove "([^"]*)" from gopath$`, iRemoveFromGopath)
	s.Step(`^I run "([^"]*)" in "([^"]*)" directory$`, iRunInDirectory)
	s.Step(`^I run "([^"]*)" in "([^"]*)"$`, iRunInDirectory)
	s.Step(`^I set "([^"]*)" to "([^"]*)"$`, iSetTo)
	s.Step(`^my workstation fulfill the requirements$`, func() error { return iRunInDirectory("./check_requirement_workstation.sh", "scripts") })
	s.Step(`^the "([^"]*)" is set to "([^"]*)"$`, theIsSetTo)
	s.Step(`^the "([^"]*)" repository exist$`, theRepositoryExist)
	s.Step(`^the directory "([^"]*)" exist$`, theDirectoryExist)
	s.Step(`^the file "([^"]*)" exist$`, theFileExist)
	s.Step(`^the output contains "([^"]*)" and "([^"]*)"$`, theOutputContainsAnd)
	s.Step(`^there is "([^"]*)" directory$`, theDirectoryExist)
	s.Step(`^there is no "([^"]*)" directory$`, thereIsNoDirectory)
}
