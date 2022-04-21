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
	"time"

	"github.com/00ahui/utils"
)

const (
	haproxyStartCmd = "if test -f $PIDFILE; then rm -f $PIDFILE; fi; nohup /usr/sbin/haproxy -Ws -f $HAPROXY_CONF -p $PIDFILE > /dev/null &"
	haproxyStopCmd  = "if test -f $PIDFILE; then PID=$(cat $PIDFILE); if test -d /proc/$PID; then /bin/kill -USR1 $PID; fi; rm -f $PIDFILE; fi"
)

func StopHAProxy(pidFile string) error {
	if status := utils.GetProcStatus(pidFile); status == utils.STATUS_UP {
		utils.LogPrintf(utils.LOG_INFO, "haproxy", "stoping haproxy '%s'", pidFile)
		cmd := fmt.Sprintf("PIDFILE=%s; %s", pidFile, haproxyStopCmd)
		if _, err := utils.ExecCommand(cmd); err != nil {
			return fmt.Errorf("stop haproxy '%s' failed: %s", pidFile, err.Error())
		}
	}
	return nil
}

func StartHAProxy(cfgFile string, pidFile string, tmout int) (string, error) {
	status := utils.STATUS_PENDING
	utils.LogPrintf(utils.LOG_INFO, "haproxy", "starting haproxy '%s'", cfgFile)
	cmd := fmt.Sprintf("HAPROXY_CONF=%s;PIDFILE=%s; %s", cfgFile, pidFile, haproxyStartCmd)
	if _, err := utils.ExecCommand(cmd); err != nil {
		return status, fmt.Errorf("starting haproxy '%s' failed: %s", cfgFile, err.Error())
	}
	for i := 0; i < tmout*100; i++ {
		if utils.FileExist(pidFile) {
			break
		}
		time.Sleep(time.Second / 100)
	}
	status = utils.GetProcStatus(pidFile)
	return status, nil
}

func RestartHAProxy(cfgFile string, pidFile string, tmout int) (string, error) {
	status := utils.STATUS_PENDING
	if err := StopHAProxy(pidFile); err != nil {
		return status, err
	}
	return StartHAProxy(cfgFile, pidFile, tmout)
}
