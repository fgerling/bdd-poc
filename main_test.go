package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	. "github.com/fgerling/bdd-poc/internal"

	"github.com/cucumber/godog"
	suse "github.com/fgerling/bdd-poc/internal/suse"
	git "gopkg.in/src-d/go-git.v4"
)

var test TestRun

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

func ReadStringAsInt(arg1 string) (int, error) {
	a, err := strconv.Atoi(arg1)
	return a, err
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

func FeatureContext(s *godog.Suite) {
	s.Step(`^I start test$`, test.IStartTest)
	s.Step(`^wait "([^"]*)"$`, wait)
	s.Step(`^"([^"]*)" exist in gopath$`, existInGopath)
	s.Step(`^I git clone "([^"]*)" into "([^"]*)"$`, iGitCloneInto)
	s.Step(`^I have "([^"]*)" in PATH$`, suse.IHaveInPATH)
	s.Step(`^I install the pattern "([^"]*)"$`, func() error { return test.IRunInDirectory("zypper -n in -t pattern SUSE-CaaSP-Management", ".") })
	s.Step(`^I remove "([^"]*)" from gopath$`, iRemoveFromGopath)
	s.Step(`^I run "([^"]*)" in "([^"]*)" directory$`, test.IRunInDirectory)
	s.Step(`^I run "([^"]*)" in "([^"]*)"$`, test.IRunInDirectory)
	s.Step(`^I set "([^"]*)" to "([^"]*)"$`, iSetTo)
	s.Step(`^my workstation fulfill the requirements$`, func() error { return test.IRunInDirectory("./check_requirement_workstation.sh", "scripts") })
	s.Step(`^the "([^"]*)" is set to "([^"]*)"$`, theIsSetTo)
	s.Step(`^the "([^"]*)" repository exist$`, theRepositoryExist)
	s.Step(`^the directory "([^"]*)" exist$`, theDirectoryExist)
	s.Step(`^the file "([^"]*)" exist$`, theFileExist)
	s.Step(`^the output contains "([^"]*)" and "([^"]*)" and "([^"]*)"$`, test.TheOutputContainsAndAnd)
	s.Step(`^the output contains "([^"]*)" and "([^"]*)"$`, test.TheOutputContainsAnd)
	s.Step(`^the output contains "([^"]*)" or "([^"]*)"$`, test.TheOutputContainsOr)
	s.Step(`^there is "([^"]*)" directory$`, theDirectoryExist)
	s.Step(`^there is no "([^"]*)" directory$`, thereIsNoDirectory)
	s.Step(`^I run VAR:"([^"]*)"$`, test.IRunVAR)
	s.Step(`^I run "([^"]*)"$`, test.Irun)
	s.Step(`^the output contains "([^"]*)"$`, test.TheOutputContains)
	s.Step(`^I have the correct go version$`, func() error { return test.IRunInDirectory("make go-version-check", "skuba") })
	s.Step(`^grep for "([^"]*)"$`, test.GrepFor)
	//--------------------Cilium-specific test functions-----------------------------------------------
	s.Step(`^I run VARS:"([^"]*)" and check for "([^"]*)" and "([^"]*)"$`, test.IRunVARSAndCheckForAnd)
	s.Step(`^VARIABLE "([^"]*)" equals ContainersFROMOutput "([^"]*)"$`, test.VARIABLEEqualsContainersFROMOutput)
	s.Step(`^VARIABLE "([^"]*)" equals ContainerFROMOutput "([^"]*)"$`, test.VARIABLEEqualsContainerFROMOutput)
	s.Step(`^I run "([^"]*)" expecting ERROR$`, test.IRunExpectingERROR)
	s.Step(`^I run VAR:"([^"]*)" expecting ERROR$`, test.IRunVARExpectingERROR)
	s.Step(`^I run VAR:"([^"]*)" expecting ERROR in VAR:"([^"]*)" directory$`, test.IRunVARExpectingERRORInVARDirectory)
	s.Step(`^the error contains "([^"]*)" and "([^"]*)"$`, test.TheErrorContainsAnd)
	s.Step(`^I run VAR:"([^"]*)" in VAR:"([^"]*)" directory$`, test.IRunVARInVARDirectory)
	s.Step(`^VARIABLES "([^"]*)" equals "([^"]*)" plus VAR:"([^"]*)"$`, test.VARIABLESEqualsPlusVAR)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)" plus VAR:"([^"]*)"$`, test.VARIABLEEqualsPlusVAR)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)" plus VAR:"([^"]*)" plus "([^"]*)"$`, test.VARIABLEEqualsPlusVARPlus)
	s.Step(`^I run "([^"]*)" in VAR:"([^"]*)" directory$`, test.IRunInVARDirectory)
	s.Step(`^VARIABLE "([^"]*)" equals "([^"]*)"$`, test.VARIABLEEquals)
	//-------------------------------------------------------------------------------------------------
	//-------------------Kured-specific test functions-------------------------------------------------
	s.Step(`^I run VARS:"([^"]*)" and IPSFromOutput$`, test.IRunVARSAndIPSFromOutput)
	s.Step(`^I run SSHCMD "([^"]*)" on MASTER$`, test.IRunSSHCMDOnMASTER)
	//-------------------Skuba Upgrade - specific test fuctions----------------------------------------
	s.Step(`^I run UPGRADE VARS:"([^"]*)" in VAR:"([^"]*)" directory$`, test.IRunUPGRADEVARSInVARDirectory)
	s.Step(`^VARIABLES "([^"]*)" equals "([^"]*)" plus Master Nodes$`, test.VARIABLESEqualsPlusMasterNodes)
	s.Step(`^VARIABLES "([^"]*)" equals "([^"]*)" plus Worker Nodes$`, test.VARIABLESEqualsPlusWorkerNodes)
	s.Step(`^VARIABLES "([^"]*)" equals "([^"]*)" plus Master Node IPS$`, test.VARIABLESEqualsPlusMasterNodeIPS)
	s.Step(`^VARIABLES "([^"]*)" equals "([^"]*)" plus Worker Node IPS$`, test.VARIABLESEqualsPlusWorkerNodeIPS)
}
