package models

import (
	"fmt"
	"myproject/pmsgo/pkg/logging"
	"myproject/pmsgo/pkg/setting"
	"myproject/pmsgo/pkg/util"
	"strings"
)

// 获取数据库中的ip的信息
func GetIpsInfo(ips string) (results []SearchPwd) {
	var (
		err    error
		ipList []string
	)

	if ips == "" {
		err = db.Table("ip_passwd").
			Select("lan_ip, user, passwds, asset_id").
			Joins("left join ip_assets on ip_assets.id=ip_passwd.asset_id").
			Scan(&results).Error
	} else {
		// ips like: 192.168.19.59,192.168.19.55
		// 处理ips的格式
		ips = strings.Trim(ips, "")
		ips = strings.Trim(ips, ",")
		ipList = strings.Split(ips, ",")
		err = db.Table("ip_passwd").
			Select("lan_ip, user, passwds, asset_id").
			Joins("left join ip_assets on ip_assets.id=ip_passwd.asset_id").
			Where("ip_assets.lan_ip in (?)", ipList).Scan(&results).Error
	}

	if err != nil {
		logging.Errorf("Fail to get all manage ips: %v", err)
	}

	return results
}

// 获取指定ip的密码
func GetIpsPwd(ips string) (results []map[string]string) {
	var (
		err       error
		pwd4Human string
		searchPwd []SearchPwd
	)
	// ips like: 192.168.19.59,192.168.19.55
	searchPwd = GetIpsInfo(ips)
	for _, pwd := range searchPwd {
		res := make(map[string]string)
		if pwd4Human, err = util.DecryptUserPwd(pwd.Passwds); err != nil {
			logging.Errorf("Fail to decrypt ip:%s passowrd: %v", pwd.LanIp, err)
			res[pwd.LanIp] = fmt.Sprintf("Fail to decrypt passowrd: %v", err)
		} else {
			res[pwd.LanIp] = pwd4Human
		}
		results = append(results, res)
	}

	return results
}

// 指定ip生成密码
func GenPwdIPs(ips string) (changeList, NotExist []string) {
	var (
		err       error
		passwd    string
		crypted   string
		ipList    []string
		searchPwd []SearchPwd
	)

	// 整理参数
	ips = strings.Trim(ips, "")
	ips = strings.Trim(ips, ",")
	ipList = strings.Split(ips, ",")

	// 查询资产表中的ip
	assetsDict := make(map[string]uint)
	searchPwd = GetIpsInfo(ips)
	for _, pwd := range searchPwd {
		assetsDict[pwd.LanIp] = pwd.AssetId
	}

	// 生成密码，并插入数据库
	for _, ip := range ipList {
		passwd = util.GeneratePasswd()
		crypted, err = util.EncryptUserPwd([]byte(passwd))
		if err != nil {
			logging.Errorf("Fail to encrypt ip: %s password:%v", ip, err)
			continue
		}

		// 检查ip是否在资产表，如在，就修改密码，否则报错
		if _, ok := assetsDict[ip]; ok {
			assetsId := assetsDict[ip]
			// 修改数据库密码
			err = db.Model(&IPPasswd{}).Where("asset_id=?", assetsId).Update("passwds", crypted).Error
			if err != nil {
				logging.Errorf("Fail to change ip: %s password: %v", ip, err)
				continue
			}
			changeList = append(changeList, fmt.Sprintf("%s:%s", ip, passwd))
		} else {
			logging.Errorf("Fail to find ip:%s in Assets", ip)
			NotExist = append(NotExist, ip)
			continue
		}
	}

	return changeList, NotExist
}

// 更新密码表的信息
func UpdateIPPasswd() {

	var (
		err      error
		ipAssets []IPAssets
	)

	logging.Info("begin to update ip_passwd")
	// 从资产表查询所有的ip
	err = db.Table("ip_assets").Find(&ipAssets).Error
	if err != nil {
		logging.Errorf("Fail to find all ips from ip_assets: %v", err)
	}
	for _, ipAsset := range ipAssets {
		item := IPPasswd{
			User: setting.SSHUser,
			//Passwds: "E2pKzD8jf5AM24xUEa81jg==",
			AssetId: ipAsset.ID,
		}
		if err := db.FirstOrCreate(&IPPasswd{}, item).Error; err != nil {
			logging.Errorf("Failed to create asset id:%s in ip_passwd, reason:%v", ipAsset.ID, err)
		}
	}

	logging.Info("end to update ip_passwd")

}
