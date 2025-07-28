package initialize

import (
	"fmt"
	"io/ioutil"
	"log"

	"go-react-admin/global"

	"gopkg.in/yaml.v2"
)

// LoadConfig 加载配置文件
func LoadConfig() {
	// 读取YAML配置文件
	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析YAML配置
	config := &global.Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	global.GlobalConfig = config

	fmt.Printf("配置文件加载成功: %+v\n", global.GlobalConfig)
	fmt.Printf("Server config: %+v\n", global.GlobalConfig.Server)
}
