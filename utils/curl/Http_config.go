package curl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// HTTPRequestConfig 用于封装HTTP请求的所有配置
type HTTPRequestConfig struct {
	Method      string
	URL         string
	Headers     map[string]string
	Payload     interface{}
	ContentType string // 默认为 "application/json"
}

var (
	apiConfigs map[string]*HTTPRequestConfig
	once       sync.Once
)

// GetAPIConfigs 返回全局的API配置映射
func GetAPIConfigs() (map[string]*HTTPRequestConfig, error) {
	once.Do(func() {
		var err error
		apiConfigs, err = loadAPIConfigs("api_configs.json")
		if err != nil {
			fmt.Println("Error loading API configs:", err)
			apiConfigs = make(map[string]*HTTPRequestConfig)
		}
	})
	return apiConfigs, nil
}

// loadAPIConfigs 加载并解析JSON配置文件
func loadAPIConfigs(filePath string) (map[string]*HTTPRequestConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var apiMap map[string][]interface{}
	err = json.Unmarshal(byteValue, &apiMap)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}
	parsedAPIConfigs := make(map[string]*HTTPRequestConfig)
	for key, value := range apiMap {
		if len(value) != 4 {
			return nil, fmt.Errorf("无效的API配置项: %s", key)
		}

		method, ok := value[0].(string)
		if !ok {
			return nil, fmt.Errorf("无效的方法类型: %s", key)
		}

		url, ok := value[1].(string)
		if !ok {
			return nil, fmt.Errorf("无效的URL类型: %s", key)
		}

		headers, ok := value[2].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("无效的Headers类型: %s", key)
		}

		contentType, ok := value[3].(string)
		if !ok {
			return nil, fmt.Errorf("无效的ContentType类型: %s", key)
		}

		headerMap := make(map[string]string)
		for k, v := range headers {
			headerMap[k], ok = v.(string)
			if !ok {
				return nil, fmt.Errorf("无效的Header值类型: %s", k)
			}
		}

		parsedAPIConfigs[key] = &HTTPRequestConfig{
			Method:      method,
			URL:         url,
			Headers:     headerMap,
			ContentType: contentType,
		}
	}

	return parsedAPIConfigs, nil
}
