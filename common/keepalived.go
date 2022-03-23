package common

import (
	"fmt"
	"log"
)

const (
	keepalivedStartCmd  = "systemctl start keepalived"
	keepalivedReloadCmd = "if test -f $PIDFILE; then kill -HUP $(cat $PIDFILE); fi"
)

func StartKeepalived(pidFile string) error {
	log.Println("starting keepalived")
	if _, err := ExecCommand(keepalivedStartCmd); err != nil {
		return err
	}
	return nil
}

func ReloadKeepalived(pidFile string) error {
	status := GetProcStatus(pidFile)
	if status != STATUS_UP {
		if err := StartKeepalived(pidFile); err != nil {
			return err
		}
	}
	log.Println("reloading keepalived")
	cmd := fmt.Sprintf("PIDFILE=%s;%s", pidFile, keepalivedReloadCmd)
	if _, err := ExecCommand(cmd); err != nil {
		return err
	}
	return nil
}
