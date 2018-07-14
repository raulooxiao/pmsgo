package v1

import (
	"github.com/gin-gonic/gin"
	"myproject/pmsgo/models"
	"myproject/pmsgo/pkg/e"
	"myproject/pmsgo/pkg/util"
	"net/http"
)

type Auth struct {
	Username string `form: "username" json: "useranme" valid:"Required; MaxSize(50)"`
	Password string `form: "password" json: "password" valid:"Required; MaxSize(50)"`
}

func GetToken(c *gin.Context) {

	var auth Auth
	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	// 参数存在
	if c.BindJSON(&auth) == nil {

		// TODO 需要增加参数校验

		isExist := models.CheckAuth(auth.Username, auth.Password)
		if isExist {
			token, err := util.GenerateToken(auth.Username, auth.Password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				code = e.SUCCESS
				data["token"] = token
			}
		} else {
			code = e.ERROR_AUTH_USER_PASSWD
		}
		// 参数不存在
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}
