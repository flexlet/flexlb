package common

import (
	"fmt"
	"log"
	"time"
)

const (
	haproxyStartCmd = "if test -f $PIDFILE; then rm -f $PIDFILE; fi; nohup /usr/sbin/haproxy -Ws -f $HAPROXY_CONF -p $PIDFILE > /dev/null &"
	haproxyStopCmd  = "if test -f $PIDFILE; then PID=$(cat $PIDFILE); if test -d /proc/$PID; then /bin/kill -USR1 $PID; fi; rm -f $PIDFILE; fi"
)

func StopHAProxy(pidFile string) error {
	if status := GetProcStatus(pidFile); status == STATUS_UP {
		log.Printf("stoping haproxy: %s\n", pidFile)
		cmd := fmt.Sprintf("PIDFILE=%s; %s", pidFile, haproxyStopCmd)
		if _, err := ExecCommand(cmd); err != nil {
			log.Printf("stop haproxy '%s' failed: %s\n", pidFile, err.Error())
			return err
		}
	}
	return nil
}

func StartHAProxy(cfgFile string, pidFile string, tmout int) (string, error) {
	status := STATUS_PENDING
	log.Printf("starting haproxy: %s\n", cfgFile)
	cmd := fmt.Sprintf("HAPROXY_CONF=%s;PIDFILE=%s; %s", cfgFile, pidFile, haproxyStartCmd)
	if _, err := ExecCommand(cmd); err != nil {
		log.Printf("starting haproxy failed: %s\n", err.Error())
		return status, err
	}
	for i := 0; i < tmout*100; i++ {
		if FileExist(pidFile) {
			break
		}
		time.Sleep(time.Second / 100)
	}
	status = GetProcStatus(pidFile)
	return status, nil
}

func RestartHAProxy(cfgFile string, pidFile string, tmout int) (string, error) {
	status := STATUS_PENDING
	if err := StopHAProxy(pidFile); err != nil {
		return status, err
	}
	return StartHAProxy(cfgFile, pidFile, tmout)
}
