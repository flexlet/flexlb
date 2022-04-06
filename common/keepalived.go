package common

import (
	"fmt"
)

const (
	keepalivedStartCmd  = "systemctl start keepalived"
	keepalivedReloadCmd = "if test -f $PIDFILE; then kill -HUP $(cat $PIDFILE); fi"
)

func StartKeepalived(pidFile string) error {
	LogPrintf(LOG_INFO, "keepalived", "starting keepalived")
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
	LogPrintf(LOG_INFO, "keepalived", "reloading keepalived")
	cmd := fmt.Sprintf("PIDFILE=%s;%s", pidFile, keepalivedReloadCmd)
	if _, err := ExecCommand(cmd); err != nil {
		return err
	}
	return nil
}
