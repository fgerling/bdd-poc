package main

import (
	"bdd-poc/cilium"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	git "gopkg.in/src-d/go-git.v4"
)

func existInGopath(arg1 string) error {
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

var Out1 []byte
var Err error

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
	//fmt.Printf("%s", fmt.Sprintf("%s", string(Out1)))
	return err
}

func theOutputContainsAnd(arg1, arg2 string) error {
	var err error
	if !strings.Contains(fmt.Sprintf("%s", string(Out1)), arg1) && strings.Contains(fmt.Sprintf("%s", string(Out1)), arg2) {
		fmt.Println("ERROR!!!")
	}
	return err
}

func grepFor(arg1 string) error {
	var err error
	tmp := strings.Split(fmt.Sprintf("%s", string(Out1)), "\n")
	for _, elem := range tmp {
		if strings.Contains(strings.ToLower(elem), arg1) {
			Out1 = []byte(elem)
		}
	}
	return err
}

func ReadStringAsInt(arg1 string) (int, error) {
	a, err := strconv.Atoi(arg1)
	return a, err
}

var VarMap map[string]string

func wait(arg1 string) error {
	temp := strings.Split(arg1, " ")
	if len(temp) > 2 || len(temp) == 1 {
		log.Println("Sorry... you've mistaken the format of time input (it's <NUMBER><1*EMPTYSPACE><WORD[seconds:minutes:hours]>")
		return nil
	} else {
		switch temp[1] {
		case "seconds":
			a, err := ReadStringAsInt(temp[0])
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
			time.Sleep(time.Duration(a) * time.Second)
		case "minutes":
			a, err := ReadStringAsInt(temp[0])
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
			time.Sleep(time.Duration(a) * time.Minute)
		case "hours":
			a, err := ReadStringAsInt(temp[0])
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
			time.Sleep(time.Duration(a) * time.Hour)
		}
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	//--------------------Cilium-specific test functions-----------------------------------------------
	s.Step(`^VARIABLE "([^"]*)" equals ContainerFROMOutput "([^"]*)"$`, cilium.VARIABLEEqualsContainerFROMOutput)
	s.Step(`^I run VAR:"([^"]*)" expecting ERROR in VAR:"([^"]*)" directory$`, cilium.IRunVARExpectingERRORInVARDirectory)
	s.Step(`^the error contains "([^"]*)" and "([^"]*)"$`, cilium.TheErrorContainsAnd)
	s.Step(`^I run VAR:"([^"]*)" in VAR:"([^"]*)" directory$`, cilium.IRunVARInVARDirectory)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)" plus VAR:"([^"]*)"$`, cilium.VARIABLEEqualsPlusVAR)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)" plus VAR:"([^"]*)" plus "([^"]*)"$`, cilium.VARIABLEEqualsPlusVARPlus)
	s.Step(`^I run "([^"]*)" in VAR:"([^"]*)" directory$`, cilium.IRunInVARDirectory)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)"$`, cilium.VARIABLEEquals)
	//-------------------------------------------------------------------------------------------------
	s.Step(`^grep for "([^"]*)"$`, grepFor)
	s.Step(`^there is "([^"]*)" directory$`, theDirectoryExsist)
	s.Step(`^I run "([^"]*)" in "([^"]*)" directory$`, iRunInDirectory)
	s.Step(`^the output contains "([^"]*)" and "([^"]*)"$`, theOutputContainsAnd)
	s.Step(`^"([^"]*)" exist in gopath$`, existInGopath)
	s.Step(`^I git clone "([^"]*)" into "([^"]*)"$`, iGitCloneInto)
	s.Step(`^wait "([^"]*)"$`, wait)
	s.Step(`^I remove "([^"]*)" from gopath$`, iRemoveFromGopath)
	s.Step(`^I run "([^"]*)" in "([^"]*)"$`, iRunIn)
	s.Step(`^I set "([^"]*)" to "([^"]*)"$`, iSetTo)
	s.Step(`^the "([^"]*)" is set to "([^"]*)"$`, theIsSetTo)
	s.Step(`^the "([^"]*)" repository exsist$`, theRepositoryExsist)
	s.Step(`^the directory "([^"]*)" exsist$`, theDirectoryExsist)
	s.Step(`^the file "([^"]*)" exsist$`, theFileExsist)
	s.Step(`^there is no "([^"]*)" directory$`, thereIsNoDirectory)
}
