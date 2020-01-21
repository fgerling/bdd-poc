package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/DATA-DOG/godog"
	git "gopkg.in/src-d/go-git.v4"
)

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

func iRunIn(command, workdir string) error {
	var splits = strings.Split(command, " ")
	cmd := exec.Command(splits[0], splits[1])
	cmd.Dir = workdir
	return cmd.Run()
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

var Out1 []byte

func iRunInDirectory(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(arg1, " ")
	cmd := exec.Command(tmp[0], tmp[1:]...)
	cmd.Dir = arg2
	Out1, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %s", err)
		return err
	}
	return err
}

func theOutputContainsAnd(arg1, arg2 string) error {
	var err error
	if !strings.Contains(fmt.Sprintf("%s", string(Out1)), arg1) && strings.Contains(fmt.Sprintf("%s", string(Out1)), arg2) {
		fmt.Println("ERROR!!!")
	}
	return err
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^there is "([^"]*)" directory$`, theDirectoryExist)
	s.Step(`^I run "([^"]*)" in "([^"]*)" directory$`, iRunInDirectory)
	s.Step(`^the output contains "([^"]*)" and "([^"]*)"$`, theOutputContainsAnd)
	s.Step(`^"([^"]*)" exist in gopath$`, existInGopath)
	s.Step(`^I git clone "([^"]*)" into "([^"]*)"$`, iGitCloneInto)
	s.Step(`^I remove "([^"]*)" from gopath$`, iRemoveFromGopath)
	s.Step(`^I run "([^"]*)" in "([^"]*)"$`, iRunIn)
	s.Step(`^I set "([^"]*)" to "([^"]*)"$`, iSetTo)
	s.Step(`^the "([^"]*)" is set to "([^"]*)"$`, theIsSetTo)
	s.Step(`^the "([^"]*)" repository exist$`, theRepositoryExist)
	s.Step(`^the directory "([^"]*)" exist$`, theDirectoryExist)
	s.Step(`^the file "([^"]*)" exist$`, theFileExist)
	s.Step(`^there is no "([^"]*)" directory$`, thereIsNoDirectory)
}
