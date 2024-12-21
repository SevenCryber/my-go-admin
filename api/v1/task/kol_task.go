package task

import (
	"encoding/json"
	"fmt"
	"github.com/SevenCryber/my-go-admin/config/task"
	"github.com/SevenCryber/my-go-admin/model/response"
	"github.com/SevenCryber/my-go-admin/utils/curl"
	"github.com/gin-gonic/gin"
)

type KolTask struct{}

// 获取当前登录用户的详情信息
func (*KolTask) List(ctx *gin.Context) {
	//var taskListConfig task.APIConfig
	taskListConfig, _ := task.GetAPIConfig("taskOrder", "OA")
	fmt.Println("taskListConfig", taskListConfig)
	postParamMap := map[string]interface{}{
		"aa": 123,
		"bb": 456,
	}
	userResponse, err := curl.Send(taskListConfig, postParamMap, "")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	// 将Body从string转成JSON对象
	var dataMap map[string]interface{}
	if unmarshalErr := json.Unmarshal([]byte(userResponse.Body), &dataMap); unmarshalErr != nil {
		// 如果解析失败，仍然可以返回原始字符串
		dataMap = map[string]interface{}{
			"raw": userResponse.Body,
		}
	}
	fmt.Print("请求体：", dataMap)

	response.NewSuccess().SetData("data", dataMap).Json(ctx)
}
