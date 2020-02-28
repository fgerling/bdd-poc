package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	cilium "github.com/fgerling/bdd-poc/features/cilium"

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
	Out1, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(Output))
	}
	cilium.Out1 = Out1
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

var Out1 []byte
var Err error
var VarMap map[string]string

func ReadStringAsInt(arg1 string) (int, error) {
	a, err := strconv.Atoi(arg1)
	return a, err
}

func VARIABLEEquals(arg1, arg2 string) error {
	var err error
	if VarMap == nil {
		VarMap = make(map[string]string)
	}
	VarMap[arg1] = arg2
	log.Printf("VAR: %s = %s\n", arg1, VarMap[arg1])
	return err
}

func wait(arg1 string) error {
	var err error
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
	return err
}

func iRunVAR(arg1 string) error {
	if VarMap["command5"] == "" {
		return Irun(cilium.VarMap[arg1])
	} else {
		return Irun(cilium.VarMap[arg1])
	}
}

func Irun(command string) error {
	return iRunInDirectory(command, ".")
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^wait "([^"]*)"$`, wait)
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
	s.Step(`^the output contains "([^"]*)" and "([^"]*)"$`, theOutputContainsAnd)
	s.Step(`^there is "([^"]*)" directory$`, theDirectoryExist)
	s.Step(`^there is no "([^"]*)" directory$`, thereIsNoDirectory)
	s.Step(`^I run VAR:"([^"]*)"$`, iRunVAR)
	s.Step(`^I run "([^"]*)"$`, Irun)
	s.Step(`^the output contains "([^"]*)"$`, theOutputContains)
	s.Step(`^I have the correct go version$`, func() error { return iRunInDirectory("make go-version-check", "skuba") })
	s.Step(`^grep for "([^"]*)"$`, grepFor)
	//--------------------Cilium-specific test functions-----------------------------------------------
	s.Step(`^VARIABLE "([^"]*)" equals ContainerFROMOutput "([^"]*)"$`, cilium.VARIABLEEqualsContainerFROMOutput)
	s.Step(`^I run "([^"]*)" expecting ERROR$`, cilium.IRunExpectingERROR)
	s.Step(`^I run VAR:"([^"]*)" expecting ERROR$`, cilium.IRunVARExpectingERROR)
	s.Step(`^I run VAR:"([^"]*)" expecting ERROR in VAR:"([^"]*)" directory$`, cilium.IRunVARExpectingERRORInVARDirectory)
	s.Step(`^the error contains "([^"]*)" and "([^"]*)"$`, cilium.TheErrorContainsAnd)
	s.Step(`^I run VAR:"([^"]*)" in VAR:"([^"]*)" directory$`, cilium.IRunVARInVARDirectory)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)" plus VAR:"([^"]*)"$`, cilium.VARIABLEEqualsPlusVAR)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)" plus VAR:"([^"]*)" plus "([^"]*)"$`, cilium.VARIABLEEqualsPlusVARPlus)
	s.Step(`^I run "([^"]*)" in VAR:"([^"]*)" directory$`, cilium.IRunInVARDirectory)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)"$`, cilium.VARIABLEEquals)
	//-------------------------------------------------------------------------------------------------
}
