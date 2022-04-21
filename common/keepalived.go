// Copyright (c) 2022 Yaohui Wang (yaohuiwang@outlook.com)
// FlexLB is licensed under Mulan PubL v2.
// You can use this software according to the terms and conditions of the Mulan PubL v2.
// You may obtain a copy of Mulan PubL v2 at:
//         http://license.coscl.org.cn/MulanPubL-2.0
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PubL v2 for more details.

package common

import (
	"fmt"

	"github.com/00ahui/utils"
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
