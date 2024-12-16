package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPResponse 包含响应体和可能的traceId
type HTTPResponse struct {
	Body    string
	TraceID string
}

// SendRequest HttpRequest 是一个通用的函数，用于发送HTTP请求（GET或POST）
func SendRequest(cfg *HTTPRequestConfig) (HTTPResponse, error) {
	var bodyReader io.Reader
	var err error

	switch cfg.ContentType {
	case "application/json":
		if cfg.Payload != nil {
			jsonData, err := json.Marshal(cfg.Payload)
			if err != nil {
				return HTTPResponse{}, err
			}
			bodyReader = bytes.NewReader(jsonData)
		}
	case "application/x-www-form-urlencoded":
		formData := url.Values{}
		if cfg.Payload != nil {
			for key, value := range cfg.Payload.(map[string]string) {
				formData.Set(key, value)
			}
			bodyReader = bytes.NewBufferString(formData.Encode())
		}
	default:
		return HTTPResponse{}, fmt.Errorf("unsupported content type: %s", cfg.ContentType)
	}

	req, err := http.NewRequest(cfg.Method, cfg.URL, bodyReader)
	if err != nil {
		return HTTPResponse{}, err
	}

	// 设置请求头
	for key, value := range cfg.Headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", cfg.ContentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return HTTPResponse{}, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return HTTPResponse{}, err
	}

	// 获取响应头中的traceId
	traceID := resp.Header.Get("trace-id") // 或者 "X-Trace-ID" 等等，取决于API的设计

	return HTTPResponse{
		Body:    string(body),
		TraceID: traceID,
	}, nil
}
