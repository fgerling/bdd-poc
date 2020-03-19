package features

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (test *TestRun) IRunInDirectory(arg1, arg2 string) error {
	var err error
	tmp := strings.Split(arg1, " ")
	cmd := exec.Command(tmp[0], tmp[1:]...)
	cmd.Dir = arg2
	//fmt.Printf("RUN: %s in %s\n", arg1, arg2)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %s", err)
		return err
	}
	test.Output = output
	test.Err = err
	//fmt.Printf("%s", fmt.Sprintf("%s", string(Out1)))
	return err
}

func (test *TestRun) IRunVARSAndIPSFromOutput(arg1 string) error {
	var err error
	for i := 0; i < 1000; i++ {
		if test.VarMap[arg1+strconv.Itoa(i)] != "" {
			err = test.IRunInDirectory(test.VarMap[arg1+strconv.Itoa(i)], ".")
			tmp1 := strings.Split(fmt.Sprintf("%s", string(test.Output)), "\n")
			for _, elem := range tmp1 {
				if strings.Contains(elem, "Node:") {
					tmp2 := strings.Split(strings.Split(strings.Replace(elem, " ", "", 100), ":")[len(strings.Split(strings.Replace(elem, " ", "", 100), ":"))-1], "/")
					if len(tmp2) == 2 {
						test.VARIABLEEquals(tmp2[0], tmp2[1])
					} else {
						fmt.Printf("Something's wrong with your kubectl describe...\n Is that even the right row? %s\n", elem)
					}
					break
				}
			}
			test.VarMap[arg1+strconv.Itoa(i)] = ""
		}
	}
	return err
}

func (test *TestRun) IRunSSHCMDOnMASTER(arg1 string) error {
	var ip string
	for key, _ := range test.VarMap {
		if test.VarMap["master-marked"] == "" {
			if strings.Contains(key, "master") /*&& strings.Contains(key, "00")*/ {
				ip = test.VarMap[key]
				test.VarMap["master-marked"] = ip
			}
		} else {
			ip = test.VarMap["master-marked"]
		}
	}
	dir, _ := os.Getwd()
	arg := append(
		[]string{"-q", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile /dev/null", "-i", filepath.Join(dir, "id_shared"),
			fmt.Sprintf("sles@%s", ip),
		},
		arg1,
	)
	cmd := exec.Command("ssh", arg...)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error! %s", err)
	}
	test.Output = output
	//fmt.Printf("%s\n", fmt.Sprintf("%s", string(test.Output)))
	return err
}
