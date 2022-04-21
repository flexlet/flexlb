package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"gitee.com/flexlb/flexlb-api/common"
	"gitee.com/flexlb/flexlb-api/models"

	"gopkg.in/yaml.v2"

	"github.com/00ahui/utils"
)

type HAProxyGlobalConfig struct {
	Daemon  bool   `yaml:"daemon"`
	MaxConn uint16 `yaml:"maxconn"`
	UID     uint16 `yaml:"uid"`
	GID     uint16 `yaml:"gid"`
	Log     string `yaml:"log"`
}

type HAProxyTimeoutConfig struct {
	HTTPRequest   string `yaml:"http-request"`
	HTTPKeepAlive string `yaml:"http-keep-alive"`
	Check         string `yaml:"check"`
	Queue         string `yaml:"queue"`
	Connect       string `yaml:"connect"`
}

type HAProxyDefaultsConfig struct {
	Mode    string               `yaml:"mode"`
	Log     string               `yaml:"log"`
	Options []string             `yaml:"options"`
	Retries uint16               `yaml:"retries"`
	MaxConn uint16               `yaml:"maxconn"`
	Timeout HAProxyTimeoutConfig `yaml:"timeout"`
}

type HAProxyConfig struct {
	ConfigDir      string                `yaml:"config_dir"`
	PidDir         string                `yaml:"pid_dir"`
	StartTimeout   int                   `yaml:"start_timeout"`
	GlobalConfig   HAProxyGlobalConfig   `yaml:"global"`
	DefaultsConfig HAProxyDefaultsConfig `yaml:"defaults"`
}

type KeepalivedGlobalDefs struct {
	RouterID             string  `yaml:"router_id"`
	VRRPSkipCheckAdvAddr bool    `yaml:"vrrp_skip_check_adv_addr"`
	VRRPGARPInterval     float32 `yaml:"vrrp_garp_interval"`
	VRRPGNAInterval      float32 `yaml:"vrrp_gna_interval"`
}

type KeepalivedConfig struct {
	ConfigFile         string               `yaml:"config_file"`
	ConfigDir          string               `yaml:"config_dir"`
	PidFile            string               `yaml:"pid_file"`
	GlobalDefs         KeepalivedGlobalDefs `yaml:"global_defs"`
	MinVirtualRouterId uint8                `yaml:"min_virtual_router_id"`
	MaxVirtualRouterId uint8                `yaml:"max_virtual_router_id"`
	AdvertInt          uint16               `yaml:"advert_int"`
	AuthType           string               `yaml:"auth_type"`
	AuthPass           string               `yaml:"auth_pass"`
}

type NodeEndpoint struct {
	IPAddress string `yaml:"ipaddress"`
	Port      uint16 `yaml:"port"`
}

type ClusterConfig struct {
	Name           string                  `yaml:"name"`
	Endpoint       string                  `yaml:"endpoint"`
	Advertize      string                  `yaml:"advertize"`
	Member         string                  `yaml:"member"`
	SecretKey      string                  `yaml:"secret_key"`
	ProbeInterval  uint16                  `yaml:"probe_interval"`
	SyncInterval   uint16                  `yaml:"sync_interval"`
	RetransmitMult uint16                  `yaml:"retransmit_mult"`
	Nodes          map[string]NodeEndpoint `yaml:"nodes,omitempty"`
}

var LB = struct {
	Name          string             `yaml:"name"`
	Host          string             `yaml:"host"`
	Port          uint16             `yaml:"port"`
	TLSHost       string             `yaml:"tls_host"`
	TLSPort       uint16             `yaml:"tls_port"`
	TLSCert       string             `yaml:"tls_cert"`
	TLSKey        string             `yaml:"tls_key"`
	TLSCACert     string             `yaml:"tls_ca_cert"`
	TLSClientCert string             `yaml:"tls_client_cert"`
	TLSClientKey  string             `yaml:"tls_client_key"`
	InstanceDir   string             `yaml:"instance_dir"`
	WatchInterval uint16             `yaml:"watch_interval"`
	LogLevel      uint8              `yaml:"log_level"`
	Cluster       ClusterConfig      `yaml:"cluster"`
	HAProxy       HAProxyConfig      `yaml:"haproxy"`
	Keepalived    KeepalivedConfig   `yaml:"keepalived"`
	Status        models.ReadyStatus `yaml:"status,omitempty"`
}{}

var (
	statusMutex    sync.Mutex
	clusterMutex   sync.Mutex
	clusterInstCfg *models.InstanceConfig
)

const (
	ReadyStatusReady    = "ready"
	ReadyStatusNotReady = "not_ready"
	ReadyStatusStarting = "starting"
)

func LoadConfig(conf string) {

	// load config file
	data, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &LB)
	if err != nil {
		log.Fatal(err)
	}

	// initialize directories
	utils.MkdirIfNotExist(LB.InstanceDir)
	utils.MkdirIfNotExist(LB.Keepalived.ConfigDir)
	utils.MkdirIfNotExist(LB.HAProxy.ConfigDir)
	utils.MkdirIfNotExist(LB.HAProxy.PidDir)

	// initialize node status map
	LB.Status = make(models.ReadyStatus)

	// initialize node endpoint map
	LB.Cluster.Nodes = map[string]NodeEndpoint{}
	LB.Cluster.Nodes[LB.Name] = NodeEndpoint{IPAddress: LB.TLSHost, Port: LB.TLSPort}

	// initialize keepalived config file
	initKeepalivedConfigFile()

	// reload keepalived
	if err := common.ReloadKeepalived(LB.Keepalived.PidFile); err != nil {
		log.Fatal(err)
	}

	// init cluster instance config
	initClusterInstCfg()
}

func UpdateNodeStatus(node string, status string) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	LB.Status[node] = status
}

func RemoveNodeStatus(node string) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	delete(LB.Status, node)
}

func UpdateNodeEndpoint(node string, nodeIp string, nodePort uint16) {
	// exist same, do nothing
	if old, exist := LB.Cluster.Nodes[node]; exist && old.IPAddress == nodeIp && old.Port == nodePort {
		return
	}

	clusterMutex.Lock()
	defer clusterMutex.Unlock()

	utils.LogPrintf(utils.LOG_DEBUG, "FlexLB", "updating cluster node endpoint '%s'", node)

	LB.Cluster.Nodes[node] = NodeEndpoint{IPAddress: nodeIp, Port: nodePort}
	UpdateClusterInstance()
}

func RemoveNodeEndpoint(node string) {
	// not exist, do nothing
	if _, exist := LB.Cluster.Nodes[node]; !exist {
		return
	}

	clusterMutex.Lock()
	defer clusterMutex.Unlock()

	utils.LogPrintf(utils.LOG_DEBUG, "FlexLB", "removing cluster node endpoint '%s'", node)

	delete(LB.Cluster.Nodes, node)
	UpdateClusterInstance()
}

func initKeepalivedConfigFile() {
	var b strings.Builder
	b.WriteString("global_defs {\n")
	b.WriteString(fmt.Sprintf("    router_id %s\n", LB.Keepalived.GlobalDefs.RouterID))
	if LB.Keepalived.GlobalDefs.VRRPSkipCheckAdvAddr {
		b.WriteString("    vrrp_skip_check_adv_addr\n")
	}
	b.WriteString(fmt.Sprintf("    vrrp_garp_interval %f\n", LB.Keepalived.GlobalDefs.VRRPGARPInterval))
	b.WriteString(fmt.Sprintf("    vrrp_gna_interval %f\n", LB.Keepalived.GlobalDefs.VRRPGNAInterval))
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("include %s/*.cfg\n", LB.Keepalived.ConfigDir))

	err := os.WriteFile(LB.Keepalived.ConfigFile, []byte(b.String()), utils.MODE_PERM_RW)
	if err != nil {
		log.Fatal(err)
	}
}

func initClusterInstCfg() {
	var (
		frontendInterface string
		frontendNetPrefix uint8
		frontendIpaddress string
		frontendPort      uint16
	)
	fields := strings.Split(LB.Cluster.Endpoint, ":")
	if len(fields) != 2 {
		log.Fatal("cluster endpoint should be host:port")
	}
	if port, err := strconv.Atoi(fields[1]); err != nil {
		log.Fatal(err)
	} else {
		frontendIpaddress = fields[0]
		frontendPort = uint16(port)
	}
	if intf, prfx, err := utils.GetInterfaceOfIP(LB.TLSHost); err != nil {
		log.Fatal(err)
	} else {
		frontendInterface = *intf
		frontendNetPrefix = uint8(*prfx)
	}

	// check-ssl with client cert/key pair
	backendDefaultServer := "inter 2s downinter 5s rise 2 fall 2 slowstart 60s maxconn 2000 maxqueue 2000 weight 100 check"

	endpoint := &models.Endpoint{
		FrontendPort:         frontendPort,
		Mode:                 "tcp",
		Balance:              "roundrobin",
		BackendDefaultServer: &backendDefaultServer,
		BackendServers:       []*models.BackendServer{},
	}

	// init cluster instance config
	clusterInstCfg = &models.InstanceConfig{
		Name:              LB.Cluster.Name,
		FrontendInterface: frontendInterface,
		FrontendNetPrefix: frontendNetPrefix,
		FrontendIpaddress: frontendIpaddress,
		Endpoints:         []*models.Endpoint{endpoint},
	}
}

// should trigger after LoadInstances
func UpdateClusterInstance() {
	// update backend servers
	clusterInstCfg.Endpoints[0].BackendServers = []*models.BackendServer{}
	for k, v := range LB.Cluster.Nodes {
		backendServer := &models.BackendServer{
			Name:      k,
			Ipaddress: v.IPAddress,
			Port:      v.Port,
		}
		clusterInstCfg.Endpoints[0].BackendServers = append(clusterInstCfg.Endpoints[0].BackendServers, backendServer)
	}

	if _, err := GetInstance(LB.Cluster.Name); err != nil {
		// not exist, create new one
		CreateInstance(clusterInstCfg)
	} else {
		// overwrite exist
		ModifyInstance(clusterInstCfg)
	}
}
