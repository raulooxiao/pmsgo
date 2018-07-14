package setting

import (
	"flag"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Cfg         *ini.File
	LogFile     string
	ExpiredTime int
	RunMode     string

	HTTPIP   string
	HTTPPort int

	ReadTimout  time.Duration
	WriteTimout time.Duration

	PageSize  int
	JwtSecret string

	SSHUser      string
	SSHPwd       string
	KeyFile      string
	KeyPwd       string
	SSHPort      int
	ConnTimeout  int
	WholeTimeout time.Duration
)

func init() {
	var (
		err     error
		cfgFile string
	)

	defaultCfg := strings.Join([]string{GetCurrentDirectory(), "conf", "app.ini"}, string(os.PathSeparator))

	flag.StringVar(&cfgFile, "f", defaultCfg, "app ini's abs path")
	flag.Parse()

	if _, err = PathExists(cfgFile); err != nil {
		log.Fatalf("%s isn't exist!", cfgFile)
	}

	Cfg, err = ini.Load(cfgFile)
	//Cfg, err = ini.Load("F:\\GOPATH\\src\\myproject\\pmsgo\\conf\\app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini: %v'", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadSSH()
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("release")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPIP = sec.Key("HTTP_IP").MustString("0.0.0.0")
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60))
	WriteTimout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60))
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
	LogFile = sec.Key("LOG_FILE").MustString("./runtime/pms-api.log")
	ExpiredTime = sec.Key("TOKEN_EXPRIE_TIME").MustInt(3)
}

func LoadSSH() {
	sec, err := Cfg.GetSection("ssh")
	if err != nil {
		log.Fatalf("Fail to get section 'ssh': %v", err)
	}

	SSHUser = sec.Key("SSH_USER").MustString("root")
	SSHPwd = sec.Key("SSH_PASSWD").MustString("root_pwd")
	KeyFile = sec.Key("KEY_FILE").MustString("/keyfile")
	KeyPwd = sec.Key("KEY_PASSWD").MustString("key_pwd")
	SSHPort = sec.Key("SSH_PORT").MustInt(22)
	ConnTimeout = sec.Key("CONNECT_TIMEOUT").MustInt(5)
	WholeTimeout = time.Duration(sec.Key("WHOLE_TIMEOUT").MustInt(600))

}
