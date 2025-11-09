package conf

import (
	"time"
)

// Bootstrap 引导配置
type Bootstrap struct {
	Server *Server `yaml:"server"`
}

// Server 服务器配置
type Server struct {
	HTTP *HTTP `yaml:"http"`
}

// HTTP HTTP服务器配置
type HTTP struct {
	Addr    string        `yaml:"addr"`
	Timeout time.Duration `yaml:"timeout"`
}
