package common

import "io/fs"

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

var (
	Debug bool = false
)
