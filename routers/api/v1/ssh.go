package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"myproject/pmsgo/models"
	"myproject/pmsgo/pkg/e"
	"myproject/pmsgo/pkg/logging"
	"myproject/pmsgo/pkg/setting"
	"myproject/pmsgo/pkg/util"
	"net/http"
	"strings"
	"time"
)

var (
	sshConfig ssh.ClientConfig
	sshClient util.SSHClient
)

type IPS struct {
	Ips string
}

func init() {
	var err error
	// 生成ssh配置
	sshConfig, err = util.CreateSSHConfig(setting.SSHUser, setting.SSHPwd, setting.KeyFile, setting.KeyPwd,
		setting.ConnTimeout)
	if err != nil {
		logging.Fatalf("Fail to Create ssh config:%v", err)
	}
}

// 批量执行
func MultiExec(execType string, ipList []string, cmd string) (succ []map[string]string, fail []map[string]string) {

	// 协程数量
	results := make(chan util.Result, 10)

	if execType == "conn" {
		for _, ip := range ipList {
			go func(_ip string) {
				results <- util.DialSuccess(_ip, setting.SSHPort, sshConfig)
			}(ip)
		}
	} else if execType == "cmd" {
		for _, ip := range ipList {
			go func(_ip, _cmd string) {
				results <- util.RunCmd(_ip, setting.SSHPort, _cmd, sshConfig)
			}(ip, cmd)
		}
	} else if execType == "chpwd" {
		for _, ips := range ipList {
			ip := strings.Split(ips, ":")[0]
			passwd := strings.Split(ips, ":")[1]
			cmd = fmt.Sprintf("echo '%s:%s' >> /tmp/xiaochpwd.log", setting.SSHUser, passwd)
			go func(_ip, _cmd string) {
				results <- util.RunCmd(_ip, setting.SSHPort, _cmd, sshConfig)
			}(ip, cmd)
		}
	}

	// 获取结果
	for i := 0; i < len(ipList); i++ {
		select {
		case res := <-results:
			data := make(map[string]string)
			data[res.Ip] = res.Msg
			if res.Rst {
				succ = append(succ, data)
			} else {
				fail = append(fail, data)
			}
		case <-time.After(setting.WholeTimeout * time.Second):
			fmt.Println("time Out")
		}
	}

	return succ, fail

}

// 修改密码
func ChangePwd(c *gin.Context) {
	var (
		code                 int
		ips                  IPS
		changeList, NotExist []string
	)

	if err := c.BindJSON(&ips); err != nil {
		logging.Errorf("解析参数异常：%v", err)
		code = e.INVALID_PARAMS
	}

	if ips.Ips == "" {
		code = e.BLANK_PARAMS
	} else {
		code = e.SUCCESS
		changeList, NotExist = models.GenPwdIPs(ips.Ips)
	}
	// 登录到机器中测试是否能联通
	succ, fail := MultiExec("chpwd", changeList, "")

	// TODO 整理结果并输出
	c.JSON(http.StatusOK, gin.H{
		"code":     code,
		"msg":      e.GetMsg(code),
		"succ":     succ,
		"fail":     fail,
		"NotExist": NotExist,
	})
}

// 查询密码
func GetPwd(c *gin.Context) {
	var (
		data []map[string]string
		code int
	)

	ips := c.DefaultQuery("ips", "")
	if ips == "" {
		code = e.BLANK_PARAMS
	} else {
		code = e.SUCCESS
		data = models.GetIpsPwd(ips)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// 检查纳管状态
func CheckStatus(c *gin.Context) {
	// ip, user
	var (
		ips    []models.SearchPwd
		ipList []string
	)

	// 从数据库中获取所有机器的列表
	ips = models.GetIpsInfo("")
	for _, _ips := range ips {
		ipList = append(ipList, _ips.LanIp)
	}

	// 登录到机器中测试是否能联通
	succ, fail := MultiExec("conn", ipList, "")

	// 整理结果并输出
	c.JSON(http.StatusOK, gin.H{
		"code":  e.SUCCESS,
		"msg":   e.GetMsg(e.SUCCESS),
		"total": len(succ) + len(fail),
		"succ":  succ,
		"fail":  fail,
	})
}

func Test(c *gin.Context) {
	//models.UpdateIPPasswd()

	jobId := 20188023
	go func(jobid int) {
		time.Sleep(1 * time.Minute)
		logging.Info("oooookkkkkkkk", ":", jobid)
	}(jobId)
	logging.Info("oooookkkkkkkk", ":", 112212)
	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"jobid":  jobId,
	})
}
