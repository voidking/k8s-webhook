package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	HTTP struct {
		Addr string // http监听端口
	}
	HTTPS struct {
		Addr   string // https监听端口
		Enable bool   // 是否启用https
		Cert   string // https证书
		Key    string // https key
	}
	// 是否开启日志Debug Level
	Debug bool
	// 时区
	TimeZone string
	// 日志路径
	LogPath string
}

var (
	instance *ServerConfig
	once     sync.Once
)

// GetConfig loads environment variables into the ServerConfig struct and returns
// a pointer to the singleton instance.

func GetConfig() *ServerConfig {
	once.Do(func() {
		// 设置读取的配置文件名和路径
		viper.SetConfigName("config")        // 配置文件名，不需要扩展名
		viper.SetConfigType("yaml")          // 配置文件类型
		viper.AddConfigPath(GetConfigPath()) // 配置文件路径

		// 读取配置
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("can not read config file: %s", err))
		}

		// 从环境变量中覆盖配置
		viper.AutomaticEnv() // 启用环境变量支持
		// 匹配环境变量时，将配置Key中的点号(.)和横杠(-)替换为下划线(_)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

		// 将配置文件反序列化为 config 对象
		var config *ServerConfig
		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
		instance = config
	})
	return instance
}

func GetConfigPath() string {
	value, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		logrus.Info("CONFIG_PATH does not exist")
	} else {
		logrus.Infof("CONFIG_PATH: %s", value)
		return filepath.Join(value)
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 向上查找直到找到go.mod文件
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			break
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // 到达根目录
			panic("There is no go.mod in project")
		}
		currentDir = parentDir
	}
	configPath := filepath.Join(currentDir, "configs")

	return configPath
}
