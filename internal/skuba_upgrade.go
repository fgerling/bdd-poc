package features

import (
	"fmt"
	"strings"
)

func (test *TestRun) VARIABLESEqualsPlusMasterNodes(arg1, arg2 string) error {
	var temp1 NodeCheck
	if test.UpgradeCheck == nil {
		test.UpgradeCheck = make(map[string]NodeCheck)
	}
	for key, _ := range test.VarMap {
		if strings.Contains(key, "master") && !strings.Contains(key, arg1) && test.UpgradeCheck[key].IP == "" {
			test.VarMap[arg1+key] = arg2 + key
			fmt.Printf("VAR: %s  = %s\n", arg1+key, arg2+key)
			temp1.IP = test.VarMap[key]
			temp1.PlanDone = false
			temp1.UPDone = false
			test.UpgradeCheck[key] = temp1
		}
	}
	return nil
}

func (test *TestRun) VARIABLESEqualsPlusWorkerNodes(arg1, arg2 string) error {
	var temp1 NodeCheck
	if test.UpgradeCheck == nil {
		test.UpgradeCheck = make(map[string]NodeCheck)
	}
	for key, _ := range test.VarMap {
		if strings.Contains(key, "worker") && !strings.Contains(key, arg1) && test.UpgradeCheck[key].IP == "" {
			test.VarMap[arg1+key] = arg2 + key
			fmt.Printf("VAR: %s  = %s\n", arg1+key, arg2+key)
			temp1.IP = test.VarMap[key]
			temp1.PlanDone = false
			temp1.UPDone = false
			test.UpgradeCheck[key] = temp1
		}
	}
	return nil
}

func (test *TestRun) VARIABLESEqualsPlusMasterNodeIPS(arg1, arg2 string) error {
	for key, _ := range test.VarMap {
		if strings.Contains(key, "master") && !strings.Contains(test.VarMap[key], "plan") && !strings.Contains(key, arg1) {
			test.VarMap[arg1+key] = arg2 + test.UpgradeCheck[key].IP
			fmt.Printf("VAR: %s  = %s\n", arg1+key, test.VarMap[arg1+key])
		}
	}
	return nil
}

func (test *TestRun) VARIABLESEqualsPlusWorkerNodeIPS(arg1, arg2 string) error {
	for key, _ := range test.VarMap {
		if strings.Contains(key, "worker") && !strings.Contains(test.VarMap[key], "plan") && !strings.Contains(key, arg1) {
			test.VarMap[arg1+key] = arg2 + test.UpgradeCheck[key].IP
			fmt.Printf("VAR: %s  = %s\n", arg1+key, test.VarMap[arg1+key])
		}
	}
	return nil
}

func (test *TestRun) IRunUPGRADEVARSInVARDirectory(arg1, arg2 string) error {
	for key, _ := range test.VarMap {
		if strings.Contains(key, arg1) && test.UpgradeCheck[key].PlanDone == false {
			//fmt.Printf("Command we run: %s\n", test.VarMap[key])
			test.Output = []byte{1}
			test.IRunInDirectory(test.VarMap[key], test.VarMap[arg2])
			temp1 := test.UpgradeCheck[key]
			temp1.PlanDone = true
			if strings.Contains(test.VarMap[key], "apply") {
				temp1.UPDone = true
			}
			test.UpgradeCheck[key] = temp1
			fmt.Println(fmt.Sprintf("%s", string(test.Output)))
			break
		}
	}
	return nil
}
