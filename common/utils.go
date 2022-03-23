package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func ExecCommand(cmd string) (int, error) {
	proc := exec.Command("bash", "-c", cmd)
	err := proc.Run()
	code := proc.ProcessState.ExitCode()
	if Debug {
		log.Printf("    %d: %s", code, cmd)
	}
	if err != nil {
		return code, err
	}
	return code, nil
}

func MkdirIfNotExist(path string) {
	d, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err1 := os.MkdirAll(path, MODE_PERM_RW)
		if err1 != nil {
			log.Fatal(err)
		}
	} else if !d.IsDir() {
		log.Fatal(fmt.Sprintf("%s already exist, but not a directory", path))
	}
}

func FileExist(f string) bool {
	_, err := os.Stat(f)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func GetProcStatus(pidFile string) string {
	if !FileExist(pidFile) {
		return STATUS_DOWN
	}
	data, err := ioutil.ReadFile(pidFile)
	if err != nil {
		log.Printf("Read file '%s' failed: %s\n", pidFile, err.Error())
		return STATUS_PENDING
	}
	line := strings.TrimSuffix(string(data), "\n")
	pid, err := strconv.ParseUint(line, 10, 32)
	if err != nil {
		log.Printf("Parse pid '%s' failed: %s\n", line, err.Error())
		return STATUS_PENDING
	}
	if FileExist(fmt.Sprintf("/proc/%d", pid)) {
		return STATUS_UP
	}
	return STATUS_DOWN
}

func GetInterfaceStatus(dev *string) string {
	cmd := fmt.Sprintf("DEV=%s;exit $(ip a s ${DEV} | grep \"${DEV}.*state UP\" | wc -l)", *dev)
	if cnt, _ := ExecCommand(cmd); cnt == 0 {
		return STATUS_DOWN
	} else {
		return STATUS_UP
	}
}

func GetIPStatus(ipaddress string, prefix *uint8, dev *string) string {
	if prefix != nil {
		ipaddress = ipaddress + fmt.Sprintf("/%d", *prefix)
	}

	var cmd string
	if dev != nil {
		cmd = fmt.Sprintf("exit $(ip a s %s | grep %s | wc -l)", *dev, ipaddress)
	} else {
		cmd = fmt.Sprintf("exit $(ip a s | grep %s | wc -l)", ipaddress)
	}

	if cnt, _ := ExecCommand(cmd); cnt == 0 {
		return STATUS_DOWN
	} else {
		return STATUS_UP
	}
}

func RandNum(max int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(max + 1)
}
