package routers

import (
	"github.com/gin-gonic/gin"
	"myproject/pmsgo/middleware/jwt"
	"myproject/pmsgo/pkg/setting"
	"myproject/pmsgo/routers/api/manage"
	"myproject/pmsgo/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	// 获取token
	r.POST("/getToken", v1.GetToken)

	// api分组
	apiv1 := r.Group("/api/v1")
	// 使用中间件
	apiv1.Use(jwt.JWTAuth())
	{
		apiv1.GET("/test", v1.Test)
		apiv1.GET("/checkstatus", v1.CheckStatus)
		apiv1.POST("/changepwd", v1.ChangePwd)
		apiv1.GET("/getpwd", v1.GetPwd)
	}

	api_manage := r.Group("/manage")
	api_manage.Use(jwt.JWTAuth())
	{
		api_manage.POST("adduser", manage.AddApiUser)
		api_manage.POST("deluser", manage.DelApiUser)
	}

	return r
}
