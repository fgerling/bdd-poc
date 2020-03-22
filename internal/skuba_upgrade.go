package features

import (
	"fmt"
	"os"
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
	var err error
	for key, _ := range test.VarMap {
		if strings.Contains(key, arg1) && test.UpgradeCheck[key].PlanDone == false {
			//fmt.Printf("Command we run: %s\n", test.VarMap[key])
			test.Output = []byte{1}
			err = test.IRunInDirectory(test.VarMap[key], test.VarMap[arg2])
			if err != nil {
				test.TreatErrors(err)
			}
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
	return err
}

func (test *TestRun) IReplaceCiliumVersionInOUTPUTAndSaveItIntoSkubaconfyamlFile() error {
	template := fmt.Sprintf("%s", string(test.Output))
	for index, row := range strings.Split(template, "\n") {
		if strings.Contains(strings.ToLower(row), "cilium") {
			version := strings.Split(template, "\n")[index+2]
			replace_version := strings.Split(version, ":")[0] + ": 1.5.1"
			template = strings.Replace(template, version, replace_version, 1)
			break
		}
	}
	f, err := os.Create("skubaconf.yaml")
	if err != nil {
		test.TreatErrors(err)
	}
	_, err = f.WriteString(template)
	if err != nil {
		test.TreatErrors(err)
	}
	f.Close()
	return err
}

func (test *TestRun) IReplaceGangwayVersionInOUTPUTAndSaveItIntoSkubaconfyamlFile() error {
	template := fmt.Sprintf("%s", string(test.Output))
	for index, row := range strings.Split(template, "\n") {
		if strings.Contains(strings.ToLower(row), "gangway") {
			version := strings.Split(template, "\n")[index+2]
			replace_version := strings.Split(version, ":")[0] + ": 2.1.0-rev4"
			template = strings.Replace(template, version, replace_version, 1)
			break
		}
	}
	f, err := os.Create("skubaconf.yaml")
	if err != nil {
		test.TreatErrors(err)
	}
	_, err = f.WriteString(template)
	if err != nil {
		test.TreatErrors(err)
	}
	f.Close()
	return err
}
