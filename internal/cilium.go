package features

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (test *TestRun) VARIABLEEqualsContainersFROMOutput(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(fmt.Sprintf("%s", string(test.Output)), "\n")
	for index, elem := range tmp {
		if strings.Contains(elem, arg2) && !strings.Contains(elem, "operator") { //--- "operator" is to exclude cilium-operator
			tmp2 := strings.Split(elem, " ")
			err = test.VARIABLEEquals(arg1+strconv.Itoa(index), tmp2[0])
		}
	}
	return err
}

func (test *TestRun) VARIABLEEqualsContainerFROMOutput(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(fmt.Sprintf("%s", string(test.Output)), "\n")
	for _, elem := range tmp {
		if strings.Contains(elem, arg2) {
			tmp2 := strings.Split(elem, " ")
			err = test.VARIABLEEquals(arg1, tmp2[0])
			break
		}
	}
	return err
}

func (test *TestRun) VARIABLEEquals(arg1, arg2 string) error {
	if test.VarMap == nil {
		test.VarMap = make(map[string]string)
	}
	test.VarMap[arg1] = arg2
	fmt.Printf(" VAR: %s = %s \n", arg1, test.VarMap[arg1])
	return nil
}

func (test *TestRun) VARIABLEEqualsPlusVARPlus(arg1, arg2, arg3, arg4 string) error {
	arg3 = test.VarMap[arg3]
	tmp := arg2 + arg3 + arg4
	err := test.VARIABLEEquals(arg1, tmp)
	return err
}

func (test *TestRun) VARIABLESEqualsPlusVAR(arg1, arg2, arg3 string) error {
	var err error
	for i := 0; i < 1000; i++ {
		if test.VarMap[arg3+strconv.Itoa(i)] != "" {
			tmp := arg2 + test.VarMap[arg3+strconv.Itoa(i)]
			err = test.VARIABLEEquals(arg1+strconv.Itoa(i), tmp)
			if err == nil {
				test.VarMap[arg3+strconv.Itoa(i)] = ""
			}
		}

	}
	return err
}

func (test *TestRun) VARIABLEEqualsPlusVAR(arg1, arg2, arg3 string) error {
	arg3 = test.VarMap[arg3]
	tmp := arg2 + arg3
	err := test.VARIABLEEquals(arg1, tmp)
	return err
}

func (test *TestRun) TheErrorContainsAnd(arg1, arg2 string) error {
	var err error
	if !strings.Contains(fmt.Sprintf("%s", test.Err), arg1) && strings.Contains(fmt.Sprintf("%s", test.Err), arg2) {
		fmt.Println("ERROR!!!")
	}
	return err
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

func (test *TestRun) IRunVARExpectingERROR(arg1 string) error {
	test.IRunExpectingERROR(test.VarMap[arg1])
	return nil
}

func (test *TestRun) IRunVARExpectingERRORInVARDirectory(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(arg1, " ")
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
	err = nil
	return err
}

func (test *TestRun) IRunVARSAndCheckForAnd(arg1, arg2, arg3 string) error {
	var err error
	for i := 0; i < 1000; i++ {
		if test.VarMap[arg1+strconv.Itoa(i)] != "" {
			err = test.IRunInDirectory(test.VarMap[arg1+strconv.Itoa(i)], ".")
			test.TheOutputContainsAnd(arg2, arg3)
			test.VarMap[arg1+strconv.Itoa(i)] = ""
		}
	}
	return err
}
