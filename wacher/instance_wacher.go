package wacher

import (
	"gitee.com/flexlb/flexlb-api/common"
	"gitee.com/flexlb/flexlb-api/config"
	"gitee.com/flexlb/flexlb-api/models"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	jobBatch uint8 = 10 // parallel batched refresh jobs
)

var (
	wacherStopped = false
	ticker        *time.Ticker   // timer
	jobGrp        sync.WaitGroup // parallel control
)

func StartInstanceWatcher() {
	log.Printf("starting instance watcher, ticker: %ds", config.Config.WatchInterval)
	ticker = time.NewTicker(time.Second * time.Duration(config.Config.WatchInterval)) // start timmer
	for {
		<-ticker.C // wait ticker

		if wacherStopped {
			break
		}

		if err := refreshKeepalived(); err != nil {
			config.ReadyStatus = models.ReadyStatusStatusKeepalivedNotReady
			continue
		}

		config.ReadyStatus = models.ReadyStatusStatusReady

		if insts := config.ListInstances(nil); len(insts) > 0 {
			if common.Debug {
				log.Printf("refreshing instances")
			}

			var jobCnt uint8 = 0
			for _, inst := range insts {
				go refreshInstance(inst) // refresh in the backend
				if jobCnt++; jobCnt >= jobBatch {
					jobGrp.Wait() // watch batch finish
					jobCnt = 0
				}
			}
			jobGrp.Wait()
		}
	}
}

func StopInstanceWacher() {
	wacherStopped = true
	if ticker != nil {
		if common.Debug {
			log.Printf("stop instance watcher")
		}
		ticker.Stop()
		ticker = nil
	}
}

func refreshKeepalived() error {
	status := common.GetProcStatus(config.Config.Keepalived.PidFile)
	if status != common.STATUS_UP {
		if err := common.StartKeepalived(config.Config.Keepalived.PidFile); err != nil {
			return err
		}
	}
	return nil
}

func refreshInstance(inst *models.Instance) {
	jobGrp.Add(1)
	cfgFile := fmt.Sprintf("%s/%s.cfg", config.Config.HAProxy.ConfigDir, inst.Config.Name)
	pidFile := fmt.Sprintf("%s/%s.pid", config.Config.HAProxy.PidDir, inst.Config.Name)
	procStatus := common.GetProcStatus(pidFile)
	vipStatus := common.GetIPStatus(inst.Config.FrontendIpaddress, &inst.Config.FrontendNetPrefix, &inst.Config.FrontendInterface)

	if vipStatus == common.STATUS_UP {
		if procStatus != common.STATUS_UP {
			// VIP启动，HAProxy未启动，则启动HAProxy
			if newStatus, err := common.StartHAProxy(cfgFile, pidFile, config.Config.HAProxy.StartTimeout); err == nil {
				config.UpdateInstanceStatus(inst.Config.Name, newStatus)
			}
		} else if inst.Status != models.InstanceStatusUp {
			// VIP启动，HAProxy已启动，但实例状态未启动，则重启HAProxy
			if newStatus, err := common.RestartHAProxy(cfgFile, pidFile, config.Config.HAProxy.StartTimeout); err == nil {
				config.UpdateInstanceStatus(inst.Config.Name, newStatus)
			}
		} else {
			// VIP启动，HAProxy已启动，实例状态已启动，则一切正常
		}
	} else {
		if procStatus == common.STATUS_UP {
			// VIP未启动，HAProxy已启动，则停止HAProxy
			common.StopHAProxy(pidFile)
		}
		config.UpdateInstanceStatus(inst.Config.Name, models.InstanceStatusDown)
	}
	jobGrp.Done()
}
