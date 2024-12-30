package task

import (
	"github.com/SevenCryber/my-go-admin/config"
	"github.com/SevenCryber/my-go-admin/model/request"
	"github.com/SevenCryber/my-go-admin/model/response"
	"github.com/SevenCryber/my-go-admin/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Order struct{}

// 获取用户列表
func (*Order) Page(ctx *gin.Context) {

	var param request.UserPage

	if err := ctx.BindQuery(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	userPages := make([]response.UserPage, 0)

	users, count := (&service.User{}).Page(param)
	if len(users) > 0 {
		for _, user := range users {
			var userPage response.UserPage
			userPage.User = user
			profile := (&service.Profile{}).GetDetailByUserId(user.Id)
			userPage.Gender = profile.Gender
			userPage.Avatar = profile.Avatar
			userPage.Address = profile.Address
			userPage.Email = profile.Email
			roleIds := (&service.UserRolesRole{}).GetRoleIdsByUserId(user.Id)
			roles := (&service.Role{}).GetListByIds(roleIds, true)
			userPage.Roles = roles
			userPages = append(userPages, userPage)
		}
	}

	response.NewSuccess().SetData("data", map[string]interface{}{
		"pageData": userPages,
		"total":    count,
	}).Json(ctx)
}

// 删除用户
func (*Order) Delete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	if id == 1 {
		response.NewError().SetMsg("不能删除根用户").Json(ctx)
		return
	}

	if id == ctx.GetInt("userId") {
		response.NewError().SetMsg("非法操作，不能删除自己").Json(ctx)
		return
	}

	if err := (&service.User{}).Delete(id); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// 添加用户
func (*Order) Add(ctx *gin.Context) {

	var param request.OrderAdd

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.Order{}).Insert(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	config.Logger.Info("数据插入完毕，准备正常返回")
	response.NewSuccess().Json(ctx)
}

// 修改用户
func (*Order) Update(ctx *gin.Context) {

	var param request.UserUpdate

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	param.Id, _ = strconv.Atoi(ctx.Param("id"))

	if err := (&service.User{}).Update(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
