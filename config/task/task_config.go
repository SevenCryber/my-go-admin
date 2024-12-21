package task

import (
	"encoding/json"
	"os"
)

// APIConfig 定义API配置结构体
type APIConfig struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Method  string `json:"method"`
	Domain  string `json:"domain"`
	Service string `json:"service"`
}

// ConfigFile 定义配置文件结构体
type ConfigFile struct {
	Apis []APIConfig `json:"apis"`
}

var apiData *ConfigFile

// APIMap 用于存储API配置的映射
var apiMap = make(map[string]map[string]APIConfig)

func InitTaskConfig() {
	file, err := os.ReadFile("config/task/task_interface.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &apiData)
	if err != nil {
		panic(err)
	}

	// 构建apiMap
	for _, api := range apiData.Apis {
		if _, exists := apiMap[api.Name]; !exists {
			apiMap[api.Name] = make(map[string]APIConfig)
		}
		apiMap[api.Name][api.Domain] = api
	}
}

// GetAPIConfig 根据name和domain获取API配置
func GetAPIConfig(name, domain string) (APIConfig, bool) {
	methods, exists := apiMap[name]
	if !exists {
		return APIConfig{}, false
	}
	api, exists := methods[domain]
	return api, exists
}
