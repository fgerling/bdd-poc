package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/cucumber/godog"
	suse "github.com/fgerling/bdd-poc/internal/suse"
	git "github.com/go-git/go-git/v5"
)

var Output []byte

//var test TestRun

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
	os.Remove(path.Join(os.Getenv("GOPATH"), "bin", file))
	return nil
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
		return errors.New(string(Output))
	}
	return err
}

func theOutputContains(arg string) error {
	if !strings.Contains(fmt.Sprintf("%s", string(Output)), arg) {
		return errors.New("Output does not contain expected argument")
	}
	return nil
}

func theOutputContainsAnd(arg1, arg2 string) error {
	if !strings.Contains(fmt.Sprintf("%s", string(Output)), arg1) && strings.Contains(fmt.Sprintf("%s", string(Output)), arg2) {
		return errors.New("Output does not contain expected arguments")
	}
	return nil
}

func theOutputContainsOr(arg1, arg2 string) error {
	if strings.Contains(fmt.Sprintf("%s", string(Output)), arg1) || strings.Contains(fmt.Sprintf("%s", string(Output)), arg2) {
		return nil
	}
	return errors.New("Output does not contain expected arguments")
}

func theOutputContainsAValidIpAddress() error {
	IP := net.ParseIP(string(Output))
	if IP == nil {
		return errors.New(fmt.Sprintf("%s is not a valid textual representation of an IP address", string(Output)))
	}
	return nil
}

func theOutputShoudMatchTheOutputTheCommand(command2 string) error {
	args := strings.Split(command2, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd2Output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(cmd2Output))
	}
	return theOutputContains(string(cmd2Output))
}

func FeatureContext(s *godog.Suite) {
	var err error
	var t TestRun
	t.RestConfig, err = GetClientRestConfig()
	if err != nil {
		log.Fatal(err)
	}

	t.ClientSet, err = t.CreateClientSet()
	if err != nil {
		log.Fatal(err)
	}
	s.Step(`^I do test "([^"]*)"$`, IDoTest)
	s.Step(`^DaemonSet "([^"]*)" in namespace "([^"]*)" exists$`, t.DaemonSetInNamespaceExists)
	s.Step(`^DaemonSet "([^"]*)" in namespace "([^"]*)" should be ready$`, t.DaemonSetInNamespaceIsReady)
	s.Step(`^Deployment "([^"]*)" in namespace "([^"]*)" exists$`, t.DeploymentInNamespaceExists)
	s.Step(`^Deployment "([^"]*)" in namespace "([^"]*)" should be ready$`, t.DeploymentInNamespaceIsReady)
	s.Step(`^ConfigMap "([^"]*)" in namespace "([^"]*)" exists$`, t.ConfigMapInNamespaceExists)
	s.Step(`^cilium ConfigMap does have the options:$`, t.ciliumConfigMapDoesHaveTheOptions)
	s.Step(`^cilium ConfigMap does not have the options:$`, t.ciliumConfigMapDoesNotHaveTheOptions)

	s.Step(`^"([^"]*)" exist in gopath$`, existInGopath)
	s.Step(`^I git clone "([^"]*)" into "([^"]*)"$`, iGitCloneInto)
	s.Step(`^I have "([^"]*)" in PATH$`, suse.IHaveInPATH)
	s.Step(`^I install the pattern "([^"]*)"$`, func() error { return iRunInDirectory("zypper -n in -t pattern SUSE-CaaSP-Management", ".") })
	s.Step(`^I remove "([^"]*)" from gopath$`, iRemoveFromGopath)
	s.Step(`^I run "([^"]*)" in "([^"]*)" directory$`, iRunInDirectory)
	s.Step(`^I run "([^"]*)" in "([^"]*)"$`, iRunInDirectory)
	s.Step(`^I set "([^"]*)" to "([^"]*)"$`, iSetTo)
	s.Step(`^my workstation fulfill the requirements$`, func() error { return iRunInDirectory("./check_requirement_workstation.sh", "scripts") })
	s.Step(`^the "([^"]*)" is set to "([^"]*)"$`, theIsSetTo)
	s.Step(`^the "([^"]*)" repository exist$`, theRepositoryExist)
	s.Step(`^the directory "([^"]*)" exist$`, theDirectoryExist)
	s.Step(`^the file "([^"]*)" exist$`, theFileExist)
	s.Step(`^the output contains "([^"]*)"$`, theOutputContains)
	s.Step(`^the output contains "([^"]*)" and "([^"]*)"$`, theOutputContainsAnd)
	s.Step(`^the output contains "([^"]*)" or "([^"]*)"$`, theOutputContainsOr)
	s.Step(`^the output contains a valid ip address$`, theOutputContainsAValidIpAddress)
	s.Step(`^the output shoud match the output the command "([^"]*)"$`, theOutputShoudMatchTheOutputTheCommand)
	s.Step(`^there is "([^"]*)" directory$`, theDirectoryExist)
	s.Step(`^there is no "([^"]*)" directory$`, thereIsNoDirectory)
	s.Step(`^I run "([^"]*)"$`, func(command string) error { return iRunInDirectory(command, ".") })
	s.Step(`^I have the correct go version$`, func() error { return iRunInDirectory("make go-version-check", "skuba") })

	//	s.BeforeScenario(func(*messages.Pickle) {
	//	})

}
