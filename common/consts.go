package common

import (
	"io/fs"
	"log"
)

const (
	MODE_PERM_RW fs.FileMode = 0600
	MODE_PERM_RO fs.FileMode = 0400
)

const (
	STATUS_UP      string = "up"
	STATUS_DOWN    string = "down"
	STATUS_PENDING string = "pending"
)

const (
	PROJECT string = "FlexLB API Server"
	VERSION string = "0.1.0"
)

const (
	LOG_DEBUG int = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
)

var (
	Logger   *log.Logger = log.Default()
	LogLevel int         = LOG_DEBUG
	LogLabel []string    = []string{"DEBUG", "INFO", "WARN", "ERROR"}
)
