package controller

import (
	"github.com/SevenCryber/my-go-admin/model/request"
	"github.com/SevenCryber/my-go-admin/model/response"
	"github.com/SevenCryber/my-go-admin/service"
	"github.com/SevenCryber/my-go-admin/utils/password"
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

	var param request.UserAdd

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if param.Username == "" {
		response.NewError().SetMsg("用户名不能为空").Json(ctx)
		return
	}

	if param.Password == "" {
		response.NewError().SetMsg("密码不能为空").Json(ctx)
		return
	}
	//response.NewError().SetMsg("用户名已存在").Json(ctx)
	user := (&service.User{}).GetDetailByUsername(param.Username)
	if user.Id > 0 {
		response.NewError().SetMsg("用户名已存在").Json(ctx)
		return
	}

	param.Password = password.Generate(param.Password)

	if err := (&service.User{}).Insert(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

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
