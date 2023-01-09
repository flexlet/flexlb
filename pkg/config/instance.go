package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/flexlet/flexlb/models"
	"github.com/flexlet/flexlb/pkg/common"
	"github.com/flexlet/utils"
	"github.com/google/go-cmp/cmp"
)

var (
	instances map[string]*models.Instance = make(map[string]*models.Instance)
	instMutex sync.Mutex
)

func LoadInstances() {
	instMutex.Lock()
	defer instMutex.Unlock()

	files, err := ioutil.ReadDir(LB.InstanceDir)
	if err != nil {
		return
	}

	for _, f := range files {
		path := fmt.Sprintf("%s/%s", LB.InstanceDir, f.Name())
		data, err := ioutil.ReadFile(path)
		if err != nil {
			utils.LogPrintf(utils.LOG_ERROR, "FlexLB", "load instance file '%s' failed: %s", path, err.Error())
			continue
		}
		var inst models.Instance
		err = json.Unmarshal(data, &inst)
		if err != nil {
			utils.LogPrintf(utils.LOG_ERROR, "FlexLB", "load instance file '%s' failed: %s", path, err.Error())
			continue
		}
		instances[inst.Config.Name] = &inst
	}
}

func UpdateInstanceStatus(node string, instName string, newStatus string) {
	instMutex.Lock()
	defer instMutex.Unlock()

	_, exist := instances[instName]
	if exist {
		instances[instName].Status[node] = newStatus
		createInstanceFile(instances[instName])
	}
}

func RemoveInstancesStatus(node string) {
	instMutex.Lock()
	defer instMutex.Unlock()
	for _, v := range instances {
		delete(v.Status, node)
		createInstanceFile(v)
	}
}

func ListInstances(fuzzyName *string) []*models.Instance {
	instMutex.Lock()
	defer instMutex.Unlock()

	matched := make([]*models.Instance, 0, len(instances))
	if fuzzyName != nil {
		for k, v := range instances {
			if strings.Contains(k, *fuzzyName) {
				inst := *v // copy instance
				matched = append(matched, &inst)
			}
		}
	} else {
		for _, v := range instances {
			inst := *v // copy instance
			matched = append(matched, &inst)
		}
	}

	return matched
}

func GetInstance(name string) (*models.Instance, error) {
	instMutex.Lock()
	defer instMutex.Unlock()

	inst, exist := instances[name]
	if !exist {
		return nil, fmt.Errorf("instance '%s' does not exist", name)
	}
	return inst, nil
}

func DeleteInstance(name string) error {
	instMutex.Lock()
	defer instMutex.Unlock()

	if _, exist := instances[name]; !exist {
		return fmt.Errorf("instance '%s' does not exist", name)
	}

	// delete keepalived config
	if err := deleteKeepalivedConfigFile(name); err != nil {
		return fmt.Errorf("delete keepalived config file failed: %s", err.Error())
	}

	// reload keepalived
	if err := common.ReloadKeepalived(LB.Keepalived.PidFile); err != nil {
		return fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// stop haproxy
	pidFile := fmt.Sprintf("%s/%s.pid", LB.HAProxy.PidDir, name)
	if err := common.StopHAProxy(pidFile); err != nil {
		return fmt.Errorf("stop haproxy failed: %s", err.Error())
	}

	cfg := instances[name].Config

	// delete haproxy config
	if err := deleteHAProxyConfigFile(cfg); err != nil {
		return fmt.Errorf("delete keepalived config file failed: %s", err.Error())
	}

	// clean instance file
	if err := deleteInstanceFile(name); err != nil {
		return fmt.Errorf("clean instance failed: %s", err.Error())
	}

	// delete instance
	delete(instances, name)

	return nil
}

func StopInstance(name string) (*models.Instance, error) {
	instMutex.Lock()
	defer instMutex.Unlock()

	inst, exist := instances[name]

	if !exist {
		return nil, fmt.Errorf("instance '%s' does not exist", name)
	}

	// delete keepalived config
	if err := deleteKeepalivedConfigFile(name); err != nil {
		return nil, fmt.Errorf("delete keepalived config file failed: %s", err.Error())
	}

	// reload keepalived
	if err := common.ReloadKeepalived(LB.Keepalived.PidFile); err != nil {
		return nil, fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// stop haproxy
	pidFile := fmt.Sprintf("%s/%s.pid", LB.HAProxy.PidDir, name)
	if err := common.StopHAProxy(pidFile); err != nil {
		return nil, fmt.Errorf("stop haproxy failed: %s", err.Error())
	}

	// delete haproxy config
	if err := deleteHAProxyConfigFile(inst.Config); err != nil {
		return nil, fmt.Errorf("delete keepalived config file failed: %s", err.Error())
	}

	// update instance status
	inst.Status[LB.Name] = utils.STATUS_DOWN

	// save instance to file
	if err := createInstanceFile(inst); err != nil {
		return nil, fmt.Errorf("save instance failed: %s", err.Error())
	}

	// notify other nodes
	GossipInstanceStatus(inst.Config.Name, utils.STATUS_DOWN)

	return inst, nil
}

func StartInstance(name string) (*models.Instance, error) {
	instMutex.Lock()
	defer instMutex.Unlock()

	inst, exist := instances[name]

	if !exist {
		return nil, fmt.Errorf("instance '%s' does not exist", name)
	}

	// create new keepalived config file
	if err := createKeepalivedConfigFile(inst.Config, inst.ID); err != nil {
		return nil, fmt.Errorf("create keepalived config file failed: %s", err.Error())
	}

	// reload keepalived
	if err := common.ReloadKeepalived(LB.Keepalived.PidFile); err != nil {
		return nil, fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// create new haproxy config file
	if err := createHAProxyConfigFile(inst.Config); err != nil {
		return nil, fmt.Errorf("create haproxy config file failed: %s", err.Error())
	}

	// update instance status
	inst.Status[LB.Name] = utils.STATUS_PENDING

	// save instance to file
	if err := createInstanceFile(inst); err != nil {
		return nil, fmt.Errorf("save instance failed: %s", err.Error())
	}

	// notify other nodes
	GossipInstanceStatus(inst.Config.Name, utils.STATUS_PENDING)

	return inst, nil
}

func ModifyInstance(cfg *models.InstanceConfig) (*models.Instance, error) {
	instMutex.Lock()
	defer instMutex.Unlock()

	inst, exist := instances[cfg.Name]
	if !exist {
		return nil, fmt.Errorf("instance '%s' does not exist", cfg.Name)
	}

	inst.Config = cfg
	inst.LastModified = time.Now().UnixMilli()

	// modify keepalived config file
	if err := createKeepalivedConfigFile(cfg, inst.ID); err != nil {
		return nil, fmt.Errorf("modify keepalived config file failed: %s", err.Error())
	}

	// reload keepalived
	if err := common.ReloadKeepalived(LB.Keepalived.PidFile); err != nil {
		return nil, fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// modify haproxy config file
	if err := createHAProxyConfigFile(cfg); err != nil {
		return nil, fmt.Errorf("modify haproxy config file failed: %s", err.Error())
	}

	// update instance status
	inst.Status[LB.Name] = utils.STATUS_PENDING

	// save instance to file
	if err := createInstanceFile(inst); err != nil {
		return nil, fmt.Errorf("save instance failed: %s", err.Error())
	}

	return inst, nil
}

func CreateInstance(cfg *models.InstanceConfig) (*models.Instance, error) {
	instMutex.Lock()
	defer instMutex.Unlock()

	if _, exist := instances[cfg.Name]; exist {
		return nil, fmt.Errorf("create instance '%s' failed, already exists", cfg.Name)
	}

	// set inst ID
	id, err := allocInstId()
	if err != nil {
		return nil, fmt.Errorf("allocate instance id failed: %s", err.Error())
	}

	var inst = &models.Instance{
		ID:           id,
		Config:       cfg,
		Status:       make(map[string]string),
		LastModified: time.Now().UnixMilli(),
	}
	inst.Status[LB.Name] = utils.STATUS_PENDING

	if err := saveInstance(inst); err != nil {
		return nil, err
	}

	instances[cfg.Name] = inst

	return inst, nil
}

func SyncInstance(other *models.Instance) error {
	instMutex.Lock()
	defer instMutex.Unlock()
	return syncInstance(other)
}

func SyncInstances(others []*models.Instance) {
	instMutex.Lock()
	defer instMutex.Unlock()
	for i := 0; i < len(others); i++ {
		other := others[i]
		syncInstance(other)
	}
}

func syncInstance(other *models.Instance) error {
	name := other.Config.Name
	exist, isExist := instances[name]
	if isExist {
		equal, err := compareInstance(exist, other)
		if err != nil {
			return err
		}
		// same config
		if equal {
			return nil
		}
		// exist is newer
		if exist.LastModified > other.LastModified {
			return nil
		}
	}

	// not exist, exist but not equal, exist is older
	if err := saveInstance(other); err != nil {
		return err
	}
	instances[name] = other
	return nil
}

func compareInstance(s *models.Instance, t *models.Instance) (bool, error) {
	// only compare id and config, not status
	if s.ID != t.ID {
		return false, nil
	}
	return cmp.Equal(s.Config, t.Config), nil
}

func saveInstance(inst *models.Instance) error {
	// create new haproxy config file
	if err := createHAProxyConfigFile(inst.Config); err != nil {
		return fmt.Errorf("create haproxy config file failed: %s", err.Error())
	}

	// create new keepalived config file
	if err := createKeepalivedConfigFile(inst.Config, inst.ID); err != nil {
		return fmt.Errorf("create keepalived config file failed: %s", err.Error())
	}

	// reload keepalived
	if err := common.ReloadKeepalived(LB.Keepalived.PidFile); err != nil {
		return fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// save instance to file
	if err := createInstanceFile(inst); err != nil {
		return fmt.Errorf("save instance failed: %s", err.Error())
	}

	return nil
}

func allocInstId() (uint8, error) {
	for i := LB.Keepalived.MinVirtualRouterId; i <= LB.Keepalived.MaxVirtualRouterId; i++ {
		exist := false
		for _, v := range instances {
			if i == v.ID {
				exist = true
				break
			}
		}
		if !exist {
			return i, nil
		}
	}
	return 0, fmt.Errorf("not enough virtual router id, allocated: %d, min: %d, max: %d",
		len(instances), LB.Keepalived.MinVirtualRouterId, LB.Keepalived.MaxVirtualRouterId)
}

func createInstanceFile(inst *models.Instance) error {
	if data, err := json.Marshal(*inst); err != nil {
		return err
	} else {
		f := fmt.Sprintf("%s/%s.json", LB.InstanceDir, inst.Config.Name)
		return ioutil.WriteFile(f, data, utils.MODE_PERM_RW)
	}
}

func deleteInstanceFile(instName string) error {
	f := fmt.Sprintf("%s/%s.json", LB.InstanceDir, instName)
	if !utils.FileExist(f) {
		return nil
	}
	return os.Remove(f)
}

func createKeepalivedConfigFile(cfg *models.InstanceConfig, id uint8) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("vrrp_instance %s {\n", cfg.Name))
	b.WriteString(fmt.Sprintf("    state %s\n    nopreempt\n", "BACKUP"))
	b.WriteString(fmt.Sprintf("    priority %d\n", utils.RandNum(100)))
	b.WriteString(fmt.Sprintf("    interface %s\n", cfg.FrontendInterface))
	b.WriteString(fmt.Sprintf("    virtual_router_id %d\n", id))
	b.WriteString(fmt.Sprintf("    advert_int %d\n", LB.Keepalived.AdvertInt))
	b.WriteString("    authentication {\n")
	b.WriteString(fmt.Sprintf("        auth_type %s\n", LB.Keepalived.AuthType))
	b.WriteString(fmt.Sprintf("        auth_pass %s\n", LB.Keepalived.AuthPass))
	b.WriteString("    }\n")
	b.WriteString("    virtual_ipaddress {\n")
	b.WriteString(fmt.Sprintf("        %s/%d\n", cfg.FrontendIpaddress, cfg.FrontendNetPrefix))
	b.WriteString("    }\n")
	b.WriteString("}\n\n")
	f := fmt.Sprintf("%s/%s.cfg", LB.Keepalived.ConfigDir, cfg.Name)
	return os.WriteFile(f, []byte(b.String()), utils.MODE_PERM_RW)
}

func deleteKeepalivedConfigFile(instName string) error {
	f := fmt.Sprintf("%s/%s.cfg", LB.Keepalived.ConfigDir, instName)
	if !utils.FileExist(f) {
		return nil
	}
	return os.Remove(f)
}

func createHAProxyConfigFile(cfg *models.InstanceConfig) error {
	var b strings.Builder

	b.WriteString("global\n")
	b.WriteString(fmt.Sprintf("    pidfile %s/%s.pid\n", LB.HAProxy.PidDir, cfg.Name))
	if LB.HAProxy.GlobalConfig.Daemon {
		b.WriteString("    daemon\n")
	}
	b.WriteString(fmt.Sprintf("    maxconn %d\n", LB.HAProxy.GlobalConfig.MaxConn))
	b.WriteString(fmt.Sprintf("    uid %d\n", LB.HAProxy.GlobalConfig.UID))
	b.WriteString(fmt.Sprintf("    gid %d\n", LB.HAProxy.GlobalConfig.GID))
	b.WriteString(fmt.Sprintf("    log %s\n", LB.HAProxy.GlobalConfig.Log))

	b.WriteString("\ndefaults\n")
	b.WriteString(fmt.Sprintf("    mode %s\n", LB.HAProxy.DefaultsConfig.Mode))
	b.WriteString(fmt.Sprintf("    log %s\n", LB.HAProxy.DefaultsConfig.Log))
	for _, o := range LB.HAProxy.DefaultsConfig.Options {
		b.WriteString(fmt.Sprintf("    option %s\n", o))
	}
	b.WriteString(fmt.Sprintf("    retries %d\n", LB.HAProxy.DefaultsConfig.Retries))
	b.WriteString(fmt.Sprintf("    timeout http-request %s\n", LB.HAProxy.DefaultsConfig.Timeout.HTTPRequest))
	b.WriteString(fmt.Sprintf("    timeout http-keep-alive %s\n", LB.HAProxy.DefaultsConfig.Timeout.HTTPKeepAlive))
	b.WriteString(fmt.Sprintf("    timeout check %s\n", LB.HAProxy.DefaultsConfig.Timeout.Check))
	b.WriteString(fmt.Sprintf("    timeout queue %s\n", LB.HAProxy.DefaultsConfig.Timeout.Queue))
	b.WriteString(fmt.Sprintf("    timeout connect %s\n", LB.HAProxy.DefaultsConfig.Timeout.Connect))
	b.WriteString(fmt.Sprintf("    maxconn %d\n", LB.HAProxy.DefaultsConfig.MaxConn))

	for _, ept := range cfg.Endpoints {
		eptName := fmt.Sprintf("%s-%d", cfg.Name, ept.FrontendPort)

		b.WriteString(fmt.Sprintf("\nfrontend %s\n", eptName))
		bind := fmt.Sprintf("%s:%d", cfg.FrontendIpaddress, ept.FrontendPort)

		if ept.FrontendOptions != nil && len(*ept.FrontendOptions) > 0 {
			bind = bind + " " + *ept.FrontendOptions
		}

		// ssl crt <inst>.pem ca-file <inst>-ca.pem verify required
		if ept.FrontendSslOptions != nil {
			frontendCrt := fmt.Sprintf("%s/%s.pem", LB.HAProxy.ConfigDir, eptName)
			utils.CreateFile(frontendCrt, ept.FrontendSslOptions.ServerCert, ept.FrontendSslOptions.ServerKey)
			bind = bind + " ssl crt " + frontendCrt
			if ept.FrontendSslOptions.CaCert != nil {
				frontendCaFile := fmt.Sprintf("%s/%s-ca.pem", LB.HAProxy.ConfigDir, eptName)
				utils.CreateFile(frontendCaFile, *ept.FrontendSslOptions.CaCert)
				bind = bind + " ca-file " + frontendCaFile
			}
			if ept.FrontendSslOptions.Verify != nil {
				bind = bind + " verify " + *ept.FrontendSslOptions.Verify
			}
		}

		b.WriteString(fmt.Sprintf("    bind %s\n", bind))

		b.WriteString(fmt.Sprintf("    default_backend %s\n", eptName))
		b.WriteString(fmt.Sprintf("    mode %s\n", ept.Mode))

		b.WriteString(fmt.Sprintf("\nbackend %s\n", eptName))
		b.WriteString(fmt.Sprintf("    mode %s\n", ept.Mode))
		for _, opt := range ept.BackendOptions {
			b.WriteString(fmt.Sprintf("    option %s\n", opt))
		}
		if ept.BackendCheckCommands != nil {
			for _, cmd := range ept.BackendCheckCommands.Commands {
				b.WriteString(fmt.Sprintf("    %s %s\n", ept.BackendCheckCommands.CheckType, cmd))
			}
		}
		if ept.BackendDefaultServer != nil {
			b.WriteString(fmt.Sprintf("    default-server %s\n", *ept.BackendDefaultServer))
		}
		b.WriteString(fmt.Sprintf("    balance %s\n", ept.Balance))
		for _, backend := range ept.BackendServers {
			server := fmt.Sprintf("%s %s:%d", backend.Name, backend.Ipaddress, backend.Port)
			if backend.Options != nil && len(*backend.Options) > 0 {
				server = server + " " + *backend.Options
			}
			// check-ssl crt <inst>-<backend>.pem ca-file <inst>-<backend>-ca.pem verify none
			if backend.CheckSslOptions != nil {
				backendCrt := fmt.Sprintf("%s/%s-%s.pem", LB.HAProxy.ConfigDir, eptName, backend.Name)
				utils.CreateFile(backendCrt, backend.CheckSslOptions.ClientCert, backend.CheckSslOptions.ClientKey)
				server = server + " check-ssl crt " + backendCrt
				if backend.CheckSslOptions.CaCert != nil {
					backendCaFile := fmt.Sprintf("%s/%s-%s-ca.pem", LB.HAProxy.ConfigDir, eptName, backend.Name)
					utils.CreateFile(backendCaFile, *backend.CheckSslOptions.CaCert)
					server = server + " ca-file " + backendCaFile
				}
				if backend.CheckSslOptions.Verify != nil {
					server = server + " verify " + *backend.CheckSslOptions.Verify
				}
			}
			b.WriteString(fmt.Sprintf("    server %s\n", server))
		}
	}

	f := fmt.Sprintf("%s/%s.cfg", LB.HAProxy.ConfigDir, cfg.Name)
	return os.WriteFile(f, []byte(b.String()), utils.MODE_PERM_RW)
}

func deleteHAProxyConfigFile(cfg *models.InstanceConfig) error {
	for _, ept := range cfg.Endpoints {
		// endpoit name
		eptName := fmt.Sprintf("%s-%d", cfg.Name, ept.FrontendPort)
		// clean frontend pem files
		if ept.FrontendSslOptions != nil {
			frontendCrt := fmt.Sprintf("%s/%s.pem", LB.HAProxy.ConfigDir, eptName)
			utils.DelFileIfExist(frontendCrt)
			if ept.FrontendSslOptions.CaCert != nil {
				frontendCaFile := fmt.Sprintf("%s/%s-ca.pem", LB.HAProxy.ConfigDir, eptName)
				utils.DelFileIfExist(frontendCaFile)
			}
		}

		// clean backend pem files
		for _, backend := range ept.BackendServers {
			if backend.CheckSslOptions != nil {
				backendCrt := fmt.Sprintf("%s/%s-%s.pem", LB.HAProxy.ConfigDir, eptName, backend.Name)
				utils.DelFileIfExist(backendCrt)
				if backend.CheckSslOptions.CaCert != nil {
					backendCaFile := fmt.Sprintf("%s/%s-%s-ca.pem", LB.HAProxy.ConfigDir, eptName, backend.Name)
					utils.DelFileIfExist(backendCaFile)
				}
			}
		}
	}
	// clean haproxy cfg
	f := fmt.Sprintf("%s/%s.cfg", LB.HAProxy.ConfigDir, cfg.Name)
	return utils.DelFileIfExist(f)
}
