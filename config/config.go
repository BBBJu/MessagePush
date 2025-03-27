package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	// 改用 mapstructure 标签（Viper 底层使用）
	AppId     string      `mapstructure:"app_id"`
	AppSecret string      `mapstructure:"app_secret"`
	MySQL     MySQLConfig `mapstructure:"mysql"`
	Redis     RedisConfig `mapstructure:"redis"`
}
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"user"`
	Pwd      string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Global variable to hold the configuration
var MyConfig AppConfig

// init function to read the configuration file
func InitConfig() {
	// 改用绝对路径直接指定配置文件（跳过路径搜索）
	viper.SetConfigFile("c:/Users/29236/Desktop/MessagePush/config/config.yaml")

	// 打印当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", dir)

	// 添加文件存在性检查
	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		log.Fatalf("Config file not found at: %s", viper.ConfigFileUsed())
	}

	if err := viper.ReadInConfig(); err != nil {
		// 更详细的错误信息
		log.Fatalf("Fatal error config file: %v \n Path used: %s", err, viper.ConfigFileUsed())
	}

	// 添加临时调试输出
	fmt.Println("Raw config content:", viper.GetString("mysql.port"))

	if err := viper.Unmarshal(&MyConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	fmt.Printf("Loaded config: %+v\n", MyConfig)
}
