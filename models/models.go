package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"myproject/pmsgo/pkg/setting"
)

// IP资产表
type IPAssets struct {
	gorm.Model
	LanIP  string `gorm:"primary_key"`
	WanIP  string
	BisId  string
	Module string
}

// ip与密码对应表
type IPPasswd struct {
	gorm.Model
	User    string `gorm:"primary_key"`
	Passwds string
	AssetId uint
}

// 操作记录表
type Recrod struct {
	ID       uint `gorm:"primary_key"`
	Exectype string
	Username string
	Params   string
	Rst      string
}

// 查询结果
type SearchPwd struct {
	LanIp string
	IPPasswd
}

// api账户表
type Pms_Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `gorm:"primary_key" json:"username"`
	Password string `json:"password"`
}

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

	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err = gorm.Open(dbType, dbConn)
	if err != nil {
		log.Fatalf("Fail to connnet 'db': %v", err)
	}

	// 不懂为啥只用单表
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// 创建表
	CreateTable()
	// 添加初始化数据
	InitData()
}

func CloseDB() {
	defer db.Close()
}

// 创建表
func CreateTable() {
	var err error
	err = db.AutoMigrate(&IPAssets{}).Error
	if err != nil {
		log.Fatalf("Fail to create table: IPAssets, %v", err)
	}

	err = db.AutoMigrate(&IPPasswd{}).Error
	if err != nil {
		log.Fatalf("Fail to create table: IPPasswd, %v", err)
	}

	err = db.AutoMigrate(&Recrod{}).Error
	if err != nil {
		log.Fatalf("Fail to create table: recrod, %v", err)
	}

	err = db.AutoMigrate(&Pms_Auth{}).Error
	if err != nil {
		log.Fatalf("Fail to create table: pms_auth, %v", err)
	}
}

// 添加初始化数据
func InitData() {
	ipList := []string{"192.168.19.59", "192.168.19.55", "192.168.19.58", "192.168.19.56", "192.168.19.235",
		"192.168.19.236", "192.168.19.237"}
	for _, ip := range ipList {
		mnip := IPAssets{
			LanIP:  ip,
			BisId:  "1001",
			Module: "pms_test",
			//Passwds: "xijfajkldafjla==1e3",
		}
		if err := db.FirstOrCreate(&IPAssets{}, mnip).Error; err != nil {
			log.Printf("Failed to create data ip:%s, reason:%v", ip, err)
		}
	}

	// pms/pms@2018
	pmsAuth := Pms_Auth{Username: "pms", Password: "a7de890663ed39b6767f2730454e17c7"}
	if err := db.FirstOrCreate(&Pms_Auth{}, pmsAuth).Error; err != nil {
		log.Fatalf("Fail to init pms_auth data: %v", err)
	}
}
