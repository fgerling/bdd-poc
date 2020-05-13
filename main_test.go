package main

import (
	"log"

	"github.com/cucumber/godog"
	suse "github.com/fgerling/bdd-poc/internal/suse"
	"github.com/fgerling/bdd-poc/internal/testrun"
)

func FeatureContext(s *godog.Suite) {
	var err error
	var t testrun.TestRun

	t.RestConfig, err = testrun.GetClientRestConfig()
	if err != nil {
		log.Fatal(err)
	}

	t.ClientSet, err = t.CreateClientSet()
	if err != nil {
		log.Fatal(err)
	}

	s.Step(`^In namespace "([^"]*)" DaemonSet "([^"]*)" exists$`, t.DaemonSetExists)
	s.Step(`^In namespace "([^"]*)" DaemonSet "([^"]*)" should be ready$`, t.DaemonSetIsReady)
	s.Step(`^In namespace "([^"]*)" Deployment "([^"]*)" exists$`, t.DeploymentExists)
	s.Step(`^In namespace "([^"]*)" Deployment "([^"]*)" should be ready$`, t.DeploymentIsReady)
	s.Step(`^In namespace "([^"]*)" ConfigMap "([^"]*)" exists$`, t.ConfigMapExists)
	s.Step(`^cilium ConfigMap does have the options:$`, t.CiliumConfigMapDoesHaveTheOptions)
	s.Step(`^cilium ConfigMap does not have the options:$`, t.CiliumConfigMapDoesNotHaveTheOptions)
	s.Step(`^In namespace "([^"]*)" no pods with labels "([^"]*)" exist$`, t.NoPodWithLabelsExist)

	s.Step(`^I run "([^"]*)" and fails$`, func(command string) error {
		return testrun.Silent(t.IRunInDirectory(command, "."))
	})

	s.Step(`^httpbin must be ready$`, t.HttpbinMustBeReady)
	s.Step(`^tblshoot must be ready$`, t.TblshootMustBeReady)

	// KUBERNETES GO-CLIENT
	s.Step(`^I exec in namespace "([^"]*)" in pod "([^"]*)" the command "([^"]*)"$`, t.IExecCommandInPod)
	s.Step(`^I exec in namespace "([^"]*)" in pod "([^"]*)" in container "([^"]*)" the command "([^"]*)"$`, t.IExecCommandInPodContainer)

	// KUBECTL
	s.Step(`^I exec with kubectl in namespace "([^"]*)" in pod "([^"]*)" the command "([^"]*)"$`, t.IKubectlExecCommandInPod)
	s.Step(`^I exec with kubectl in namespace "([^"]*)" in pod "([^"]*)" in container "([^"]*)" the command "([^"]*)"$`, t.IKubectlExecCommandInPodContainer)

	s.Step(`^I apply the manifest "([^"]*)"$`, t.IApplyManifest)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, t.ISendRequestTo)
	s.Step(`^I send "([^"]*)" request to "([^"]*)" and fails$`, func(method, path string) error {
		return testrun.Silent(t.ISendRequestTo(method, path))
	})

	s.Step(`^there is no resource "([^"]*)" in "([^"]*)" namespace$`, t.ThereIsNoResourceInNamespace)

	s.Step(`^I resolve "([^"]*)"$`, t.IResolve)
	s.Step(`^I resolve "([^"]*)" and fails$`, func(fqdn string) error {
		return testrun.Silent(t.IResolve(fqdn))
	})

	s.Step(`^I reverse resolve "([^"]*)"$`, t.IReverseResolve)
	s.Step(`^I reverse resolve "([^"]*)" and fails$`, func(ip string) error {
		return testrun.Silent(t.IReverseResolve(ip))
	})

	s.Step(`^"([^"]*)" exist in gopath$`, t.ExistInGopath)
	s.Step(`^I git clone "([^"]*)" into "([^"]*)"$`, t.IGitCloneInto)
	s.Step(`^I have "([^"]*)" in PATH$`, suse.IHaveInPATH)
	s.Step(`^I install the pattern "([^"]*)"$`, func() error { return t.IRunInDirectory("zypper -n in -t pattern SUSE-CaaSP-Management", ".") })
	s.Step(`^I remove "([^"]*)" from gopath$`, t.IRemoveFromGopath)
	s.Step(`^I set "([^"]*)" to "([^"]*)"$`, t.ISetTo)
	s.Step(`^my workstation fulfill the requirements$`, func() error { return t.IRunInDirectory("./check_requirement_workstation.sh", "scripts") })
	s.Step(`^the "([^"]*)" is set to "([^"]*)"$`, t.TheIsSetTo)
	s.Step(`^the "([^"]*)" repository exist$`, t.TheRepositoryExist)
	s.Step(`^the directory "([^"]*)" exist$`, t.TheDirectoryExist)
	s.Step(`^the file "([^"]*)" exist$`, t.TheFileExist)
	s.Step(`^the output contains "([^"]*)"$`, t.TheOutputContains)
	s.Step(`^the output contains "([^"]*)" and "([^"]*)"$`, t.TheOutputContainsAnd)
	s.Step(`^the output contains "([^"]*)" or "([^"]*)"$`, t.TheOutputContainsOr)
	s.Step(`^the output contains a valid ip address$`, t.TheOutputContainsAValidIpAddress)
	s.Step(`^the output shoud match the output the command "([^"]*)"$`, t.TheOutputShoudMatchTheOutputTheCommand)
	s.Step(`^there is "([^"]*)" directory$`, t.TheDirectoryExist)
	s.Step(`^there is no "([^"]*)" directory$`, t.ThereIsNoDirectory)
	s.Step(`^I run "([^"]*)" in "([^"]*)" directory$`, t.IRunInDirectory)
	s.Step(`^I run "([^"]*)"$`, func(command string) error { return t.IRunInDirectory(command, ".") })
	s.Step(`^I have the correct go version$`, func() error { return t.IRunInDirectory("make go-version-check", "skuba") })
}
