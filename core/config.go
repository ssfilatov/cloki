package cloki

import (
	//"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type ClickhouseConfig struct {
	URL            string `yaml:"url"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Database       string `yaml:"database"`
	LogTableName   string `yaml:"table"`
	MaxConnections int    `yaml:"max_connections"`
}

type ServerConfig struct {
	HTTPListenHost string `yaml:"http_listen_host"`
	HTTPListenPort int    `yaml:"http_listen_port"`

	HTTPServerReadTimeout  time.Duration `yaml:"http_server_read_timeout"`
	HTTPServerWriteTimeout time.Duration `yaml:"http_server_write_timeout"`
	HTTPServerIdleTimeout  time.Duration `yaml:"http_server_idle_timeout"`
}

type Config struct {
	Debug      bool              `yaml:"debug"`
	LabelList  *[]string         `yaml:"label_list"`
	Server     *ServerConfig     `yaml:"server"`
	Clickhouse *ClickhouseConfig `yaml:"clickhouse"`
}

//func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
//	f.BoolVar(&cfg.Debug, "debug", false, "Use debug logging.")
//	f.StringVar(&cfg.Server.HTTPListenHost, "server.http-listen-host", "", "HTTP server listen host.")
//	f.IntVar(&cfg.Server.HTTPListenPort, "server.http-listen-port", 80, "HTTP server listen port.")
//	f.DurationVar(&cfg.Server.HTTPServerReadTimeout, "server.http-read-timeout", 30*time.Second, "Read timeout for HTTP server")
//	f.DurationVar(&cfg.Server.HTTPServerWriteTimeout, "server.http-write-timeout", 30*time.Second, "Write timeout for HTTP server")
//	f.DurationVar(&cfg.Server.HTTPServerIdleTimeout, "server.http-idle-timeout", 120*time.Second, "Idle timeout for HTTP server")
//}

func LoadConfig(filename string, cfg interface{}) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading config file: %s", err)
	}

	return yaml.UnmarshalStrict(buf, cfg)
}
