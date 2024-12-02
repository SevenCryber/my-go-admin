package middleware

import (
	"github.com/SevenCryber/my-go-admin/initialize/message"
	"github.com/SevenCryber/my-go-admin/service"
	"github.com/SevenCryber/my-go-admin/token"
	"github.com/gin-gonic/gin"
)

// 鉴权中间件
func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userClaims, err := token.ParseToken(ctx)
		if err != nil {
			message.Error(ctx, 401, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("userId", userClaims.UserId)
		ctx.Set("username", userClaims.Username)
		ctx.Set("roleCode", userClaims.CurrentRoleCode)

		// 超级管理员不做鉴权
		if userClaims.CurrentRoleCode == "SUPER_ADMIN" {
			ctx.Next()
			return
		}

		permission := (&service.Permission{}).GetDetailByPathAndMethod(ctx.FullPath(), "")
		if permission.Id <= 0 {
			ctx.Next()
			return
		}

		if role := (&service.Role{}).GetDetailByCode(userClaims.CurrentRoleCode); role.Id > 0 {
			if !(&service.RolePermissionsPermission{}).CheckHasPermission(role.Id, permission.Id) {
				message.Error(ctx, "您目前暂无此权限，请联系管理员申请权限")
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
