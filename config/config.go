package config

import (
	"flexlb/common"
	"flexlb/models"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
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
	RouterID             string `yaml:"router_id"`
	VRRPSkipCheckAdvAddr bool   `yaml:"vrrp_skip_check_adv_addr"`
	VRRPGARPInterval     uint16 `yaml:"vrrp_garp_interval"`
	VRRPGNAInterval      uint16 `yaml:"vrrp_gna_interval"`
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

type GlobalConfig struct {
	ConfigDir     string           `yaml:"config_dir"`
	WatchInterval uint8            `yaml:"watch_interval"`
	HAProxy       HAProxyConfig    `yaml:"haproxy"`
	Keepalived    KeepalivedConfig `yaml:"keepalived"`
}

var Config GlobalConfig
var ReadyStatus string

func LoadConfig(f string) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatal(err)
	}
	common.MkdirIfNotExist(Config.ConfigDir)
	common.MkdirIfNotExist(Config.Keepalived.ConfigDir)
	common.MkdirIfNotExist(Config.HAProxy.ConfigDir)
	common.MkdirIfNotExist(Config.HAProxy.PidDir)

	initKeepalivedConfigFile()

	if err := common.ReloadKeepalived(Config.Keepalived.PidFile); err != nil {
		log.Fatal(err)
	}

	ReadyStatus = models.ReadyStatusStatusReady
}

func initKeepalivedConfigFile() {
	var b strings.Builder
	b.WriteString("global_defs {\n")
	b.WriteString(fmt.Sprintf("    router_id %s\n", Config.Keepalived.GlobalDefs.RouterID))
	if Config.Keepalived.GlobalDefs.VRRPSkipCheckAdvAddr {
		b.WriteString("    vrrp_skip_check_adv_addr\n")
	}
	b.WriteString(fmt.Sprintf("    vrrp_garp_interval %d\n", Config.Keepalived.GlobalDefs.VRRPGARPInterval))
	b.WriteString(fmt.Sprintf("    vrrp_gna_interval %d\n", Config.Keepalived.GlobalDefs.VRRPGNAInterval))
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("include %s/*.cfg\n", Config.Keepalived.ConfigDir))

	err := os.WriteFile(Config.Keepalived.ConfigFile, []byte(b.String()), common.MODE_PERM_RW)
	if err != nil {
		log.Fatal(err)
	}
}
