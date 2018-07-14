package models

import (
	"errors"
	"myproject/pmsgo/pkg/setting"
	"myproject/pmsgo/pkg/util"
)

// 检查用户密码是否一致
func CheckAuth(username, password string) bool {
	var auth Pms_Auth
	hash_pwd := util.EncryptPass(username, password, setting.JwtSecret)
	db.Select("id").Where(Pms_Auth{Username: username, Password: hash_pwd}).First(&auth)
	if auth.ID > 0 {
		return true
	}

	return false
}

// 新建api调用账户
func AddUser(auth Pms_Auth) error {
	var user Pms_Auth
	var err error

	// 检查账户是否存在
	db.Where(Pms_Auth{Username: auth.Username}).First(&user)
	if user.ID > 0 {
		return errors.New("该账户已存在!")
	}

	// 创建账户
	auth.Password = util.EncryptPass(auth.Username, auth.Password, setting.JwtSecret)
	if err = db.FirstOrCreate(&Pms_Auth{}, auth).Error; err != nil {
		return err
	}

	return nil
}

// 删除api调用账户
func DelUser(auth Pms_Auth) error {
	if err := db.Delete(&Pms_Auth{}, auth).Error; err != nil {
		return err
	}

	return nil
}
