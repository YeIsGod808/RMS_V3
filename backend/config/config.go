package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var globalConfig = new(GlobalConfig)

type GlobalConfig struct {
	*SvrConfig   `mapstructure:"svr_config"`
	*LogConfig   `mapstructure:"log" json:"log" yaml:"log"`
	*Neo4jConfig `mapstructure:"neo4j"`
	*DbConfig    `mapstructure:"mysql"`
	*JwtConfig   `mapstructure:"jwt"`
}

type SvrConfig struct {
	Name string `mapstructure:"name"` // 服务name
	Host string `mapstructure:"host"` // 服务host
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // gin模式
}

type DbConfig struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Database    string `mapstructure:"database"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`  // 最大空闲连接数
	MaxOpenConn int    ` mapstructure:"max_open_conn"` // 最大打开的连接数
	MaxIdleTime int64  ` mapstructure:"max_idle_time"` // 连接最大空闲时间
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"file_name"`
	LogPath    string `mapstructure:"log_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	LoggerBuf  string `mapstructure:"logger_buf"`
}
type Neo4jConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type JwtConfig struct {
	Issuer  string `json:"issuer"`
	JwtSalt string `json:"jwt_salt"`
}

func Init() (err error) {
	// 自动推导项目根目录
	configFile := GetRootDir() + "/config/config.yaml"
	viper.SetConfigFile(configFile)

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()

	if err != nil {
		// 读取配置信息错误
		fmt.Printf("viper.ReadInConfig() failed: %v\n", err)
		return fmt.Errorf("viper.ReadInConfig() failed: %v\n", err)
	}

	if err = viper.Unmarshal(globalConfig); err != nil {
		fmt.Printf("viper.Unmarshal failed: %v\n", err)
		return fmt.Errorf("viper.Unmarshal failed: %v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改...")
		if err = viper.Unmarshal(globalConfig); err != nil {
			fmt.Printf("viper.Unmarshal failed: %v\n", err)
		}
	})
	return nil
}

func GetGlobalConfig() *GlobalConfig {
	return globalConfig
}
