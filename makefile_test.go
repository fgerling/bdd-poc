package main

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	git "gopkg.in/src-d/go-git.v4"
	"os"
	"os/exec"
	"path"
	"strings"
)

func exsistInGopath(arg1 string) error {
	return theFileExsist(path.Join(os.Getenv("GOPATH"), "bin"))
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

func theRepositoryExsist(repository string) error {
	return theFileExsist(path.Join(repository, ".git/"))
}

func theDirectoryExsist(dir string) error {
	return theFileExsist(dir)
}

func theFileExsist(file string) error {
	_, err := os.Stat(file)
	return err
}

func thereIsNoDirectory(target string) error {
	return os.RemoveAll(target)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^"([^"]*)" exsist in gopath$`, exsistInGopath)
	s.Step(`^I git clone "([^"]*)" into "([^"]*)"$`, iGitCloneInto)
	s.Step(`^I remove "([^"]*)" from gopath$`, iRemoveFromGopath)
	s.Step(`^I run "([^"]*)" in "([^"]*)"$`, iRunIn)
	s.Step(`^I set "([^"]*)" to "([^"]*)"$`, iSetTo)
	s.Step(`^the "([^"]*)" is set to "([^"]*)"$`, theIsSetTo)
	s.Step(`^the "([^"]*)" repository exsist$`, theRepositoryExsist)
	s.Step(`^the directory "([^"]*)" exsist$`, theDirectoryExsist)
	s.Step(`^the file "([^"]*)" exsist$`, theFileExsist)
	s.Step(`^there is no "([^"]*)" directory$`, thereIsNoDirectory)
}
