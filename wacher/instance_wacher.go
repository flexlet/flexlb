package wacher

import (
	"fmt"
	"sync"
	"time"

	"gitee.com/flexlb/flexlb-api/common"
	"gitee.com/flexlb/flexlb-api/config"
	"gitee.com/flexlb/flexlb-api/models"
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
	common.LogPrintf(common.LOG_INFO, "FlexLB", "starting instance watcher, ticker: %ds", config.LB.WatchInterval)
	ticker = time.NewTicker(time.Second * time.Duration(config.LB.WatchInterval)) // start timmer
	for {
		<-ticker.C // wait ticker

		if wacherStopped {
			break
		}

		if err := reconcileKeepalived(); err != nil {
			common.LogPrintf(common.LOG_ERROR, "FlexLB", "reconcile keepalvied failed")
			config.UpdateNodeStatus(config.LB.Name, config.ReadyStatusNotReady)
			config.GossipNodeStatus()
			continue
		}

		config.UpdateNodeStatus(config.LB.Name, config.ReadyStatusReady)
		config.GossipNodeStatus()

		if insts := config.ListInstances(nil); len(insts) > 0 {
			common.LogPrintf(common.LOG_DEBUG, "FlexLB", "reconcile all instances")

			var jobCnt uint8 = 0
			for _, inst := range insts {
				go reconcileInstance(inst) // refresh in the backend
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
		common.LogPrintf(common.LOG_DEBUG, "FlexLB", "stop instance watcher")
		ticker.Stop()
		ticker = nil
	}
}

func reconcileKeepalived() error {
	status := common.GetProcStatus(config.LB.Keepalived.PidFile)
	// start keepalived if not up
	if status != common.STATUS_UP {
		if err := common.StartKeepalived(config.LB.Keepalived.PidFile); err != nil {
			return err
		}
	}
	return nil
}

func reconcileInstance(inst *models.Instance) {
	jobGrp.Add(1)
	cfgFile := fmt.Sprintf("%s/%s.cfg", config.LB.HAProxy.ConfigDir, inst.Config.Name)
	pidFile := fmt.Sprintf("%s/%s.pid", config.LB.HAProxy.PidDir, inst.Config.Name)
	procStatus := common.GetProcStatus(pidFile)
	vipStatus := common.GetIPStatus(inst.Config.FrontendIpaddress, &inst.Config.FrontendNetPrefix, &inst.Config.FrontendInterface)

	if vipStatus == common.STATUS_UP {
		if procStatus != common.STATUS_UP {
			// vip up, haproxy not up, then start haproxy
			if newStatus, err := common.StartHAProxy(cfgFile, pidFile, config.LB.HAProxy.StartTimeout); err == nil {
				config.UpdateInstanceStatus(config.LB.Name, inst.Config.Name, newStatus)
				// notify other nodes
				config.GossipInstanceStatus(inst.Config.Name, newStatus)
			}
		} else if inst.Status[config.LB.Name] != common.STATUS_UP {
			// vip up, haproxy up, but inst down, then restart haproxy
			if newStatus, err := common.RestartHAProxy(cfgFile, pidFile, config.LB.HAProxy.StartTimeout); err == nil {
				config.UpdateInstanceStatus(config.LB.Name, inst.Config.Name, newStatus)
				// notify other nodes
				config.GossipInstanceStatus(inst.Config.Name, newStatus)
			}
		} else {
			// vip up, haproxy up, inst up, then do nothing
			// notify other nodes
			config.GossipInstanceStatus(inst.Config.Name, common.STATUS_UP)
		}
	} else {
		if procStatus == common.STATUS_UP {
			// vip down, haproxy up, then stop haproxy
			common.StopHAProxy(pidFile)
		}
		// inst is down, update status
		config.UpdateInstanceStatus(config.LB.Name, inst.Config.Name, common.STATUS_DOWN)
		// notify other nodes
		config.GossipInstanceStatus(inst.Config.Name, common.STATUS_DOWN)
	}
	jobGrp.Done()
}
