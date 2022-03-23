package config

import (
	"encoding/json"
	"flexlb/common"
	"flexlb/models"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	instances map[string]*models.Instance = make(map[string]*models.Instance)
	instMutex sync.Mutex
)

func LoadInstances() {
	instMutex.Lock()
	defer instMutex.Unlock()

	files, err := ioutil.ReadDir(Config.ConfigDir)
	if err != nil {
		return
	}

	for _, f := range files {
		path := fmt.Sprintf("%s/%s", Config.ConfigDir, f.Name())
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("load instance file '%s' failed: %s", path, err.Error())
			continue
		}
		var inst models.Instance
		err = json.Unmarshal(data, &inst)
		if err != nil {
			log.Printf("load instance file '%s' failed: %s", path, err.Error())
			continue
		}
		instances[inst.Config.Name] = &inst
	}
}

func UpdateInstanceStatus(instName string, newStatus string) {
	instMutex.Lock()
	defer instMutex.Unlock()

	_, exist := instances[instName]
	if exist {
		instances[instName].Status = newStatus
		createInstanceFile(instances[instName])
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

	// clean instance file
	if err := deleteInstanceFile(name); err != nil {
		return fmt.Errorf("clean instance failed: %s", err.Error())
	}

	// delete instance
	delete(instances, name)

	// delete keepalived config
	if err := deleteKeepalivedConfigFile(name); err != nil {
		return fmt.Errorf("delete keepalived config file failed: %s", err.Error())
	}

	// reload keepalived
	if err := common.ReloadKeepalived(Config.Keepalived.PidFile); err != nil {
		return fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// stop haproxy
	pidFile := fmt.Sprintf("%s/%s.pid", Config.HAProxy.PidDir, name)
	if err := common.StopHAProxy(pidFile); err != nil {
		return fmt.Errorf("stop haproxy failed: %s", err.Error())
	}

	// delete haproxy config
	if err := deleteHAProxyConfigFile(name); err != nil {
		return fmt.Errorf("delete keepalived config file failed: %s", err.Error())
	}

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
	if err := common.ReloadKeepalived(Config.Keepalived.PidFile); err != nil {
		return nil, fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// stop haproxy
	pidFile := fmt.Sprintf("%s/%s.pid", Config.HAProxy.PidDir, name)
	if err := common.StopHAProxy(pidFile); err != nil {
		return nil, fmt.Errorf("stop haproxy failed: %s", err.Error())
	}

	// delete haproxy config
	if err := deleteHAProxyConfigFile(name); err != nil {
		return nil, fmt.Errorf("delete keepalived config file failed: %s", err.Error())
	}

	// update instance status
	inst.Status = models.InstanceStatusDown

	// save instance to file
	if err := createInstanceFile(inst); err != nil {
		return nil, fmt.Errorf("save instance failed: %s", err.Error())
	}

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
	if err := common.ReloadKeepalived(Config.Keepalived.PidFile); err != nil {
		return nil, fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// create new haproxy config file
	if err := createHAProxyConfigFile(inst.Config); err != nil {
		return nil, fmt.Errorf("create haproxy config file failed: %s", err.Error())
	}

	// update instance status
	inst.Status = models.InstanceStatusPending

	// save instance to file
	if err := createInstanceFile(inst); err != nil {
		return nil, fmt.Errorf("save instance failed: %s", err.Error())
	}

	return inst, nil
}

func ModifyInstance(name string, cfg *models.InstanceConfig) (*models.Instance, error) {
	instMutex.Lock()
	defer instMutex.Unlock()

	inst, exist := instances[name]
	if !exist {
		return nil, fmt.Errorf("instance '%s' does not exist", name)
	}
	if name != cfg.Name {
		return nil, fmt.Errorf("does not support to modify instance name '%s' to '%s'", name, cfg.Name)
	}

	inst.Config = cfg

	// modify keepalived config file
	if err := createKeepalivedConfigFile(cfg, inst.ID); err != nil {
		return nil, fmt.Errorf("modify keepalived config file failed: %s", err.Error())
	}

	// reload keepalived
	if err := common.ReloadKeepalived(Config.Keepalived.PidFile); err != nil {
		return nil, fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// modify haproxy config file
	if err := createHAProxyConfigFile(cfg); err != nil {
		return nil, fmt.Errorf("modify haproxy config file failed: %s", err.Error())
	}

	// update instance status
	inst.Status = models.InstanceStatusPending

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
	id, err := newInstId()
	if err != nil {
		return nil, fmt.Errorf("allocate instance id failed: %s", err.Error())
	}

	var inst models.Instance

	inst.ID = id
	inst.Config = cfg
	inst.Status = models.InstanceStatusPending

	// create new haproxy config file
	if err := createHAProxyConfigFile(cfg); err != nil {
		return nil, fmt.Errorf("create haproxy config file failed: %s", err.Error())
	}

	// create new keepalived config file
	if err := createKeepalivedConfigFile(cfg, id); err != nil {
		return nil, fmt.Errorf("create keepalived config file failed: %s", err.Error())
	}

	// reload keepalived
	if err := common.ReloadKeepalived(Config.Keepalived.PidFile); err != nil {
		return nil, fmt.Errorf("reload keepalived failed: %s", err.Error())
	}

	// save instance to file
	if err := createInstanceFile(&inst); err != nil {
		return nil, fmt.Errorf("save instance failed: %s", err.Error())
	}

	instances[cfg.Name] = &inst

	return &inst, nil
}

func newInstId() (uint8, error) {
	for i := Config.Keepalived.MinVirtualRouterId; i <= Config.Keepalived.MaxVirtualRouterId; i++ {
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
		len(instances), Config.Keepalived.MinVirtualRouterId, Config.Keepalived.MaxVirtualRouterId)
}

func createInstanceFile(inst *models.Instance) error {
	if data, err := json.Marshal(*inst); err != nil {
		return err
	} else {
		f := fmt.Sprintf("%s/%s.json", Config.ConfigDir, inst.Config.Name)
		return ioutil.WriteFile(f, data, common.MODE_PERM_RW)
	}
}

func deleteInstanceFile(instName string) error {
	f := fmt.Sprintf("%s/%s.json", Config.ConfigDir, instName)
	if !common.FileExist(f) {
		return nil
	}
	return os.Remove(f)
}

func createKeepalivedConfigFile(cfg *models.InstanceConfig, id uint8) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("vrrp_instance %s {\n", cfg.Name))
	b.WriteString(fmt.Sprintf("    state %s\n", cfg.State))
	b.WriteString(fmt.Sprintf("    priority %d\n", cfg.Priority))
	b.WriteString(fmt.Sprintf("    interface %s\n", cfg.FrontendInterface))
	b.WriteString(fmt.Sprintf("    virtual_router_id %d\n", id))
	b.WriteString(fmt.Sprintf("    advert_int %d\n", Config.Keepalived.AdvertInt))
	b.WriteString("    authentication {\n")
	b.WriteString(fmt.Sprintf("        auth_type %s\n", Config.Keepalived.AuthType))
	b.WriteString(fmt.Sprintf("        auth_pass %s\n", Config.Keepalived.AuthPass))
	b.WriteString("    }\n")
	b.WriteString("    virtual_ipaddress {\n")
	b.WriteString(fmt.Sprintf("        %s/%d\n", cfg.FrontendIpaddress, cfg.FrontendNetPrefix))
	b.WriteString("    }\n")
	b.WriteString("}\n\n")
	f := fmt.Sprintf("%s/%s.cfg", Config.Keepalived.ConfigDir, cfg.Name)
	return os.WriteFile(f, []byte(b.String()), common.MODE_PERM_RW)
}

func deleteKeepalivedConfigFile(instName string) error {
	f := fmt.Sprintf("%s/%s.cfg", Config.Keepalived.ConfigDir, instName)
	if !common.FileExist(f) {
		return nil
	}
	return os.Remove(f)
}

func createHAProxyConfigFile(cfg *models.InstanceConfig) error {
	var b strings.Builder

	b.WriteString("global\n")
	b.WriteString(fmt.Sprintf("    pidfile %s/%s.pid\n", Config.HAProxy.PidDir, cfg.Name))
	if Config.HAProxy.GlobalConfig.Daemon {
		b.WriteString("    daemon\n")
	}
	b.WriteString(fmt.Sprintf("    maxconn %d\n", Config.HAProxy.GlobalConfig.MaxConn))
	b.WriteString(fmt.Sprintf("    uid %d\n", Config.HAProxy.GlobalConfig.UID))
	b.WriteString(fmt.Sprintf("    gid %d\n", Config.HAProxy.GlobalConfig.GID))
	b.WriteString(fmt.Sprintf("    log %s\n", Config.HAProxy.GlobalConfig.Log))

	b.WriteString("\ndefaults\n")
	b.WriteString(fmt.Sprintf("    mode %s\n", Config.HAProxy.DefaultsConfig.Mode))
	b.WriteString(fmt.Sprintf("    log %s\n", Config.HAProxy.DefaultsConfig.Log))
	for _, o := range Config.HAProxy.DefaultsConfig.Options {
		b.WriteString(fmt.Sprintf("    option %s\n", o))
	}
	b.WriteString(fmt.Sprintf("    retries %d\n", Config.HAProxy.DefaultsConfig.Retries))
	b.WriteString(fmt.Sprintf("    timeout http-request %s\n", Config.HAProxy.DefaultsConfig.Timeout.HTTPRequest))
	b.WriteString(fmt.Sprintf("    timeout http-keep-alive %s\n", Config.HAProxy.DefaultsConfig.Timeout.HTTPKeepAlive))
	b.WriteString(fmt.Sprintf("    timeout check %s\n", Config.HAProxy.DefaultsConfig.Timeout.Check))
	b.WriteString(fmt.Sprintf("    timeout queue %s\n", Config.HAProxy.DefaultsConfig.Timeout.Queue))
	b.WriteString(fmt.Sprintf("    timeout connect %s\n", Config.HAProxy.DefaultsConfig.Timeout.Connect))
	b.WriteString(fmt.Sprintf("    maxconn %d\n", Config.HAProxy.DefaultsConfig.MaxConn))

	b.WriteString(fmt.Sprintf("\nfrontend %s\n", cfg.Name))
	b.WriteString(fmt.Sprintf("    bind %s:%d %s\n", cfg.FrontendIpaddress, cfg.FrontendPort, *cfg.FrontendOptions))
	b.WriteString(fmt.Sprintf("    default_backend %s\n", cfg.Name))
	b.WriteString(fmt.Sprintf("    mode %s\n", cfg.Mode))

	b.WriteString(fmt.Sprintf("\nbackend %s\n", cfg.Name))
	b.WriteString(fmt.Sprintf("    mode %s\n", cfg.Mode))
	for _, opt := range cfg.BackendOptions {
		b.WriteString(fmt.Sprintf("    option %s\n", opt))
	}
	if cfg.BackendCheckCommands != nil {
		for _, cmd := range cfg.BackendCheckCommands.Commands {
			b.WriteString(fmt.Sprintf("    %s %s\n", cfg.BackendCheckCommands.CheckType, cmd))
		}
	}
	b.WriteString(fmt.Sprintf("    default-server %s\n", *cfg.BackendDefaultServer))
	b.WriteString(fmt.Sprintf("    balance %s\n", cfg.Balance))
	for _, svr := range cfg.BackendServers {
		b.WriteString(fmt.Sprintf("    server %s %s:%d %s\n", svr.Name, svr.Ipaddress, svr.Port, *svr.Options))
	}

	f := fmt.Sprintf("%s/%s.cfg", Config.HAProxy.ConfigDir, cfg.Name)
	return os.WriteFile(f, []byte(b.String()), common.MODE_PERM_RW)
}

func deleteHAProxyConfigFile(instName string) error {
	f := fmt.Sprintf("%s/%s.cfg", Config.HAProxy.ConfigDir, instName)
	if !common.FileExist(f) {
		return nil
	}
	return os.Remove(f)
}
