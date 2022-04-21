// Copyright (c) 2022 Yaohui Wang (yaohuiwang@outlook.com)
// FlexLB is licensed under Mulan PubL v2.
// You can use this software according to the terms and conditions of the Mulan PubL v2.
// You may obtain a copy of Mulan PubL v2 at:
//         http://license.coscl.org.cn/MulanPubL-2.0
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PubL v2 for more details.

package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitee.com/flexlb/flexlb-api/models"
	"github.com/hashicorp/memberlist"
	"github.com/00ahui/utils"
)

const (
	MsgTypeCreateInstance byte = 'c'
	MsgTypeModifyInstance byte = 'm'
	MsgTypeDeleteInstance byte = 'd'
	MsgTypeStopInstance   byte = 'p'
	MsgTypeStartInstance  byte = 'r'
	MsgTypeInstanceStatus byte = 's'
	MsgTypeNodeStatus     byte = 'S'
	MsgTypeNodeEndpoint   byte = 'e'
)

var gossip *utils.Gossip

func StartClusterGossip() error {

	// split gossip endpoint to host & port
	field := strings.Split(LB.Cluster.Advertize, ":")
	if len(field) != 2 {
		return fmt.Errorf("gossip endpoint format error: %s, should be <host>:<port>", LB.Cluster.Advertize)
	}
	bindAddr := field[0]
	bindPort, err := strconv.Atoi(field[1])
	if err != nil {
		return fmt.Errorf("gossip endpoint format error: %s, should be <host>:<port>", LB.Cluster.Advertize)
	}

	gossip = utils.GossipWith(LB.Name, bindAddr, bindPort)
	gossip.SecretKey = []byte(LB.Cluster.SecretKey)
	gossip.ProbeInterval = int(LB.Cluster.ProbeInterval)
	gossip.SyncInterval = int(LB.Cluster.SyncInterval)
	gossip.RetransmitMult = int(LB.Cluster.RetransmitMult)
	gossip.NotifyMsgHandler = notifyMsgHandler
	gossip.LocalStateHandler = localStateHandler
	gossip.MergeRemoteStateHandler = mergeRemoteStateHandler
	gossip.NotifyJoinHandler = notifyJoinHandler
	gossip.NotifyLeaveHandler = notifyLeaveHandler
	if err := gossip.Start(&LB.Cluster.Member); err != nil {
		return err
	}

	// set ready state
	UpdateNodeStatus(LB.Name, ReadyStatusReady)

	// notify node ready
	GossipNodeStatus()

	// notify node endpoint
	GossipNodeEndpoint()

	return nil
}

func GossipCreateInstance(inst *models.Instance) error {
	data, err := json.Marshal(*inst)
	if err != nil {
		return err
	}
	gossip.Broadcast(append([]byte{MsgTypeCreateInstance}, data...))
	return nil
}

func GossipModifyInstance(inst *models.Instance) error {
	data, err := json.Marshal(*inst)
	if err != nil {
		return err
	}
	gossip.Broadcast(append([]byte{MsgTypeModifyInstance}, data...))
	return nil
}

func GossipDeleteInstance(inst string) {
	data := []byte(inst)
	gossip.Broadcast(append([]byte{MsgTypeDeleteInstance}, data...))
}

func GossipStopInstance(inst string) {
	data := []byte(inst)
	gossip.Broadcast(append([]byte{MsgTypeStopInstance}, data...))
}

func GossipStartInstance(inst string) {
	data := []byte(inst)
	gossip.Broadcast(append([]byte{MsgTypeStartInstance}, data...))
}

func GossipInstanceStatus(inst string, status string) {
	data := []byte(fmt.Sprintf("%s:%s:%s", LB.Name, inst, status))
	gossip.Broadcast(append([]byte{MsgTypeInstanceStatus}, data...))
}

func GossipNodeStatus() {
	data := []byte(LB.Name + ":" + LB.Status[LB.Name])
	gossip.Broadcast(append([]byte{MsgTypeNodeStatus}, data...))
}

func GossipNodeEndpoint() {
	data := []byte(fmt.Sprintf("%s:%s:%d", LB.Name, LB.TLSHost, LB.TLSPort))
	gossip.Broadcast(append([]byte{MsgTypeNodeEndpoint}, data...))
}

func notifyMsgHandler(msg []byte) {
	msgType := msg[0]
	msgData := msg[1:]
	switch msgType {
	// create or modify instance
	case MsgTypeCreateInstance, MsgTypeModifyInstance:
		var inst = &models.Instance{}
		if err := json.Unmarshal(msgData, inst); err != nil {
			return
		}
		SyncInstance(inst)
	// delete instance
	case MsgTypeDeleteInstance:
		name := string(msgData)
		DeleteInstance(name)
	// stop instance
	case MsgTypeStopInstance:
		name := string(msgData)
		StopInstance(name)
	// start instance
	case MsgTypeStartInstance:
		name := string(msgData)
		StartInstance(name)
	// ready status
	case MsgTypeInstanceStatus:
		// node:inst:status
		v := strings.Split(string(msgData), ":")
		if len(v) == 3 {
			UpdateInstanceStatus(v[0], v[1], v[2])
		}
	// ready status
	case MsgTypeNodeStatus:
		utils.LogPrintf(utils.LOG_DEBUG, "FlexLB", "received node status msg: '%s'", string(msgData))
		// node:status
		v := strings.Split(string(msgData), ":")
		if len(v) == 2 {
			UpdateNodeStatus(v[0], v[1])
		}
	// node endpoint
	case MsgTypeNodeEndpoint:
		utils.LogPrintf(utils.LOG_DEBUG, "FlexLB", "received node endpoint msg: '%s'", string(msgData))
		// node:ip:port
		v := strings.Split(string(msgData), ":")
		if len(v) == 3 {
			if port, err := strconv.Atoi(v[2]); err == nil {
				UpdateNodeEndpoint(v[0], v[1], uint16(port))
			}
		}
	}
}

// send local instances
func localStateHandler() []byte {
	insts := ListInstances(nil)
	if data, err := json.Marshal(insts); err != nil {
		return []byte{}
	} else {
		return data
	}
}

// receive remote instances
func mergeRemoteStateHandler(data []byte) {
	if len(data) > 0 {
		var insts = []models.Instance{}
		if err := json.Unmarshal(data, &insts); err != nil {
			return
		}
		// copy pointers
		var pinsts = []*models.Instance{}
		for i := 0; i < len(insts); i++ {
			pinsts = append(pinsts, &insts[i])
		}
		SyncInstances(pinsts)
	}
}

// node leave
func notifyLeaveHandler(node *memberlist.Node) {
	utils.LogPrintf(utils.LOG_INFO, "FlexLB", "node '%s' left cluster", node.Name)
	RemoveNodeStatus(node.Name)
	RemoveNodeEndpoint(node.Name)
	RemoveInstancesStatus(node.Name)
}

// node join
func notifyJoinHandler(node *memberlist.Node) {
	utils.LogPrintf(utils.LOG_INFO, "FlexLB", "node '%s' joined cluster", node.Name)
	if node.Name != LB.Name {
		UpdateNodeStatus(node.Name, ReadyStatusStarting)
		if gossip.Ready {
			utils.LogPrintf(utils.LOG_DEBUG, "FlexLB", "node '%s' joined cluster, gossip local node", node.Name)
			// gossip local node to new joined node
			GossipNodeStatus()
			GossipNodeEndpoint()
		}
	}
}
