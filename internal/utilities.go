package features

import (
	"errors"
	"fmt"
	"strings"
)

func (test *TestRun) IStartTest() error {
	test.Output = []byte{1}
	test.VarMap = make(map[string]string)
	test.Err = nil
	return nil
}

func (test *TestRun) TheOutputContains(arg string) error {
	if !strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg) {
		return errors.New("Output does not contain expected argument")
	}
	return nil
}

func (test *TestRun) TheOutputContainsAnd(arg1, arg2 string) error {
	if !strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg1) && strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg2) {
		return errors.New("Output does not contain expected arguments")
	}
	return nil
}

func (test *TestRun) TheOutputContainsOr(arg1, arg2 string) error {
	if !strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg1) || strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg2) {
		return errors.New("Output does not contain expected arguments")
	}
	return nil
}

func (test *TestRun) IRunVAR(arg1 string) error {
	if test.VarMap["command5"] == "" {
		return test.Irun(test.VarMap[arg1])
	} else {
		return test.Irun(test.VarMap[arg1])
	}
}

func (test *TestRun) Irun(command string) error {
	return test.IRunInDirectory(command, ".")
}

func (test *TestRun) GrepFor(arg1 string) error {
	var err error
	tmp := strings.Split(fmt.Sprintf("%s", string(test.Output)), "\n")
	for _, elem := range tmp {
		if strings.Contains(strings.ToLower(elem), arg1) {
			test.Output = []byte(elem)
		}
	}
	return err
}
