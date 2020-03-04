package cilium

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var Out1 []byte
var Err error
var VarMap map[string]string

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

func IRunVARInVARDirectory(arg1, arg2 string) error {
	arg1 = VarMap[arg1]
	err := IRunInVARDirectory(arg1, arg2)
	return err
}

func IRunInVARDirectory(arg1, arg2 string) error {
	arg2 = VarMap[arg2]
	err := iRunInDirectory(arg1, arg2)
	return err
}

func VARIABLEEqualsContainersFROMOutput(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(fmt.Sprintf("%s", string(Out1)), "\n")
	for index, elem := range tmp {
		if strings.Contains(elem, arg2) {
			tmp2 := strings.Split(elem, " ")
			err = VARIABLEEquals(arg1+strconv.Itoa(index), tmp2[0])
		}
	}
	return err
}

func VARIABLEEqualsContainerFROMOutput(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(fmt.Sprintf("%s", string(Out1)), "\n")
	for _, elem := range tmp {
		if strings.Contains(elem, arg2) {
			tmp2 := strings.Split(elem, " ")
			err = VARIABLEEquals(arg1, tmp2[0])
			break
		}
	}
	return err
}

func VARIABLEEquals(arg1, arg2 string) error {
	var err error
	if VarMap == nil {
		VarMap = make(map[string]string)
	}
	VarMap[arg1] = arg2
	fmt.Printf(" VAR: %s = %s \n", arg1, VarMap[arg1])
	return err
}

func VARIABLEEqualsPlusVARPlus(arg1, arg2, arg3, arg4 string) error {
	arg3 = VarMap[arg3]
	tmp := arg2 + arg3 + arg4
	err := VARIABLEEquals(arg1, tmp)
	return err
}

func VARIABLESEqualsPlusVAR(arg1, arg2, arg3 string) error {
	var err error
	for i := 0; i < 1000; i++ {
		if VarMap[arg3+strconv.Itoa(i)] != "" {
			tmp := arg2 + VarMap[arg3+strconv.Itoa(i)]
			err = VARIABLEEquals(arg1+strconv.Itoa(i), tmp)
			if err == nil {
				VarMap[arg3+strconv.Itoa(i)] = ""
			}
		}

	}
	return err
}

func VARIABLEEqualsPlusVAR(arg1, arg2, arg3 string) error {
	arg3 = VarMap[arg3]
	tmp := arg2 + arg3
	err := VARIABLEEquals(arg1, tmp)
	return err
}

func TheErrorContainsAnd(arg1, arg2 string) error {
	var err error
	if !strings.Contains(fmt.Sprintf("%s", Err), arg1) && strings.Contains(fmt.Sprintf("%s", Err), arg2) {
		fmt.Println("ERROR!!!")
	}
	return err
}

func IRunExpectingERROR(arg1 string) error {
	var err error
	tmp := strings.Split(arg1, " ")
	cmd := exec.Command(tmp[0], tmp[1:]...)
	Out1, err = cmd.CombinedOutput()
	if err != nil {
		if !strings.Contains(fmt.Sprintf("%s", err), "exit") && strings.Contains(fmt.Sprintf("%s", err), "28") {
			fmt.Fprintf(os.Stdout, "error: %s", err)
			return err
		}
	}
	Err = err
	err = nil
	return err
}

func IRunVARExpectingERROR(arg1 string) error {
	IRunExpectingERROR(VarMap[arg1])
	return nil
}

func IRunVARExpectingERRORInVARDirectory(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(arg1, " ")
	cmd := exec.Command(tmp[0], tmp[1:]...)
	cmd.Dir = arg2
	Out1, err = cmd.CombinedOutput()
	if err != nil {
		if !strings.Contains(fmt.Sprintf("%s", err), "exit code") && strings.Contains(fmt.Sprintf("%s", err), "28") {
			fmt.Fprintf(os.Stdout, "error: %s", err)
			return err
		}
	}
	Err = err
	err = nil
	return err
}

func theOutputContainsAnd(arg1, arg2 string) error {
	if !strings.Contains(fmt.Sprintf("%s", string(Out1)), arg1) && strings.Contains(fmt.Sprintf("%s", string(Out1)), arg2) {
		return errors.New("Output does not contain expected arguments")
	}
	return nil
}

func IRunVARSAndCheckForAnd(arg1, arg2, arg3 string) error {
	var err error
	for i := 0; i < 1000; i++ {
		if VarMap[arg1+strconv.Itoa(i)] != "" {
			err = iRunInDirectory(VarMap[arg1+strconv.Itoa(i)], ".")
			theOutputContainsAnd(arg2, arg3)
			VarMap[arg1+strconv.Itoa(i)] = ""
		}
	}
	return err
}