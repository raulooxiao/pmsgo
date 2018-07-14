package main

import (
	"fmt"
)

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"log"
	"myproject/pmsgo/models"
	"myproject/pmsgo/pkg/setting"
	"strings"
)

var db *gorm.DB

func init() {
	var (
		err                                  error
		dbType, dbName, user, password, host string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()

	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", user, password, host, dbName)
	db, err = gorm.Open(dbType, dbConn)
	if err != nil {
		log.Fatalf("Fail to connnet 'db': %v", err)
	}

	// 不懂为啥只用单表
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()

}

type Pms_Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Pms_Auth
	db.Exec("select * from pms_auth where username=xiao and password=AES_DECRYPT('xiao123', 'sjsops')").
		First(&auth)
	fmt.Println(auth.ID, auth.Username, auth.Password)
	if auth.ID > 0 {
		return true
	}

	return false
}

type SearchPwd struct {
	LanIp string
	models.IPPasswd
}

func GetIpsPwd(ips string) []map[string]string {
	var searchpwd []SearchPwd
	var ipList []string
	var results []map[string]string
	// ips like: 192.168.19.59,192.168.19.55
	// 处理ips的格式
	ips = strings.Trim(ips, "")
	ips = strings.Trim(ips, ",")
	ipList = strings.Split(ips, ",")
	err := db.Table("ip_passwd").
		Select("lan_ip, user, passwds, asset_id").
		Joins("left join ip_assets on ip_assets.id=ip_passwd.asset_id").
		Where("ip_assets.lan_ip in (?)", ipList).Scan(&searchpwd).Error
	if err != nil {
		fmt.Printf("Fail to get ips:%s passwd: %v", ips, err)
	}
	for _, pwd := range searchpwd {
		res := make(map[string]string)
		res[pwd.LanIp] = pwd.Passwds
		results = append(results, res)
	}

	return results
}

func main() {
	ips := "192.168.19.239,192.168.19.56"
	GetIpsPwd(ips)
	var err error
	err = db.Model(&models.IPPasswd{}).Where("asset_id=?", 20).Update("passwds", "1xiaoisgoodman").Error
	if err != nil {
		log.Printf("Fail to change ip: %s password: %v", "xiaofdafd", err)
	}

}
