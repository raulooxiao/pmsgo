package manage

import (
	"github.com/gin-gonic/gin"
	"myproject/pmsgo/models"
	"myproject/pmsgo/pkg/e"
	"myproject/pmsgo/pkg/logging"
	"net/http"
)

// 增加api调用者
func AddApiUser(c *gin.Context) {
	var (
		err  error
		auth models.Pms_Auth
	)

	code := e.SUCCESS

	if err = c.BindJSON(&auth); err != nil {
		logging.Errorf("解析参数异常：%v", err)
		code = e.INVALID_PARAMS
	}

	if err = models.AddUser(auth); err != nil {
		logging.Errorf("添加Api调用账户失败，异常信息：%v", err)
		code = e.FAIL_CREATE_USER
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}

// 删除api调用者
func DelApiUser(c *gin.Context) {
	var (
		err  error
		auth models.Pms_Auth
	)
	code := e.SUCCESS

	if err = c.BindJSON(&auth); err != nil {
		logging.Errorf("解析参数异常：%v", err)
		code = e.INVALID_PARAMS
	}

	err = models.DelUser(auth)
	if err != nil {
		logging.Errorf("删除Api调用账户失败，错误信息：%v", err)
		code = e.FAIL_DELETE_USER
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}
