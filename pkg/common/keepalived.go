package common

import (
	"fmt"

	"github.com/flexlet/utils"
)

const (
	keepalivedStartCmd  = "systemctl start keepalived"
	keepalivedReloadCmd = "if test -f $PIDFILE; then kill -HUP $(cat $PIDFILE); fi"
)

func StartKeepalived(pidFile string) error {
	utils.LogPrintf(utils.LOG_INFO, "keepalived", "starting keepalived")
	if _, err := utils.ExecCommand(keepalivedStartCmd); err != nil {
		return err
	}
	return nil
}

func ReloadKeepalived(pidFile string) error {
	status := utils.GetProcStatus(pidFile)
	if status != utils.STATUS_UP {
		if err := StartKeepalived(pidFile); err != nil {
			return err
		}
	}
	utils.LogPrintf(utils.LOG_INFO, "keepalived", "reloading keepalived")
	cmd := fmt.Sprintf("PIDFILE=%s;%s", pidFile, keepalivedReloadCmd)
	if _, err := utils.ExecCommand(cmd); err != nil {
		return err
	}
	return nil
}
