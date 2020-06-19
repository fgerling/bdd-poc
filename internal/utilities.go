package features

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func (test *TestRun) TheOutputContainsAndAnd(arg1, arg2, arg3 string) error {
	if !strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg1) && strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg2) && strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg3) {
		return errors.New("Output does not contain expected arguments")
	}
	return nil
}

func (test *TestRun) TheOutputContainsOr(arg1, arg2 string) error {
	//fmt.Printf("text:\n%s\nvars: %s   %s\n ", strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg1, arg2)
	if strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg1) || strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg2) {
	} else {
		return errors.New("Output does not contain expected arguments")
	}
	return nil
}

func (test *TestRun) TheOutputContainsOrOr(arg1, arg2, arg3 string) error {
	//fmt.Printf("text:\n%s\nvars: %s   %s\n ", strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg1, arg2)
	if strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg1) || strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg2) || strings.Contains(strings.ToLower(fmt.Sprintf("%s", string(test.Output))), arg3) {
	} else {
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

func (test *TestRun) IRunVARInVARDirectory(arg1, arg2 string) error {
	arg1 = test.VarMap[arg1]
	err := test.IRunInVARDirectory(arg1, arg2)
	fmt.Println(fmt.Sprintf("%s", string(test.Output)))
	return err
}

func (test *TestRun) IRunInVARDirectory(arg1, arg2 string) error {
	arg2 = test.VarMap[arg2]
	err := test.IRunInDirectory(arg1, arg2)
	fmt.Println(fmt.Sprintf("%s", string(test.Output)))
	return err
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

func (test *TestRun) TreatErrors(err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nError: %v\n", err)
	}
}

func (test *TestRun) VARIABLEIGetFromCONFIG(arg1 string) error {
	if test.VarMap == nil {
		test.IStartTest()
	}
	var config Config
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		test.TreatErrors(err)
	}
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		test.TreatErrors(err)
	}
	test.Config = config
	test.VarMap[arg1] = config.ClusterDir
	fmt.Printf("VAR: %s = %s", arg1, test.VarMap[arg1])
	return nil
}

func (test *TestRun) IRunExpectingERRORInVARDirectory(arg1, arg2, arg3 string) error {
	var err error
	tmp := strings.Split(arg1, " ")
	cmd := exec.Command(tmp[0], tmp[1:]...)
	cmd.Dir = arg3
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR: %v\n grep: %s", err, arg2)
		fmt.Printf("OUTPUT: %s\n", fmt.Sprintf("%s", string(output)))
		if !strings.Contains(fmt.Sprintf("%s", err), arg2) {
			fmt.Fprintf(os.Stdout, "error: %s", err)
			return nil
		}
	} else {
		test.Err = err
	}
	test.Output = output
	return nil
}

func (test *TestRun) IRunVARExpectingERROR(arg1 string) error {
	test.IRunExpectingERROR(test.VarMap[arg1])
	return nil
}

func (test *TestRun) IRunVARExpectingERRORInVARDirectory(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(test.VarMap[arg1], " ")
	cmd := exec.Command(tmp[0], tmp[1:]...)
	cmd.Dir = arg2
	output, err := cmd.CombinedOutput()
	if err != nil {
		if !strings.Contains(fmt.Sprintf("%s", err), "exit code") && strings.Contains(fmt.Sprintf("%s", err), "28") {
			fmt.Fprintf(os.Stdout, "error: %s", err)
			return err
		}
	}
	test.Output = output
	test.Err = err
	return nil
}

func (test *TestRun) IRunExpectingERROR(arg1 string) error {
	var err error
	tmp := strings.Split(arg1, " ")
	cmd := exec.Command(tmp[0], tmp[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if !strings.Contains(fmt.Sprintf("%s", err), "exit") && strings.Contains(fmt.Sprintf("%s", err), "28") {
			fmt.Fprintf(os.Stdout, "error: %s", err)
			return err
		}
	}
	test.Output = output
	test.Err = err
	err = nil
	return err
}
