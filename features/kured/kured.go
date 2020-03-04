package kured

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fgerling/bdd-poc/features/cilium"
)

var Out1 []byte

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

func IRunVARSAndIPSFromOutput(arg1 string) error {
	var err error
	for i := 0; i < 1000; i++ {
		if cilium.VarMap[arg1+strconv.Itoa(i)] != "" {
			err = iRunInDirectory(cilium.VarMap[arg1+strconv.Itoa(i)], ".")
			tmp1 := strings.Split(fmt.Sprintf("%s", string(Out1)), "\n")
			for _, elem := range tmp1 {
				if strings.Contains(elem, "Node:") {
					tmp2 := strings.Split(strings.Split(strings.Replace(elem, " ", "", 100), ":")[len(strings.Split(strings.Replace(elem, " ", "", 100), ":"))-1], "/")
					if len(tmp2) == 2 {
						cilium.VARIABLEEquals(tmp2[0], tmp2[1])
					} else {
						fmt.Printf("Something's wrong with your kubectl describe...\n Is that even the right row? %s\n", elem)
					}
					break
				}
			}
			cilium.VarMap[arg1+strconv.Itoa(i)] = ""
		}
	}
	return err
}

func IRunSSHCMDOnMASTER(arg1 string) error {
	var ip string
	for key, _ := range cilium.VarMap {
		if strings.Contains(key, "master") /*&& strings.Contains(key, "00")*/ {
			ip = cilium.VarMap[key]
		}
	}
	arg := append(
		[]string{"-q", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile /dev/null", "-i", "id_shared",
			fmt.Sprintf("sles@%s", ip),
		},
		arg1,
	)
	Out1, err := exec.Command("ssh", arg...).CombinedOutput()
	if err != nil {
		log.Printf("Error! %s", err)
	}
	fmt.Printf("%s\n", fmt.Sprintf("%s", string(Out1)))
	return err
}
