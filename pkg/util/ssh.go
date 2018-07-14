package util

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type Result struct {
	Rst  bool
	Msg  string
	Code int
	Ip   string
}

// ssh client
type SSHClient struct {
	Ip            string
	User          string
	Password      string
	PrivateKey    string
	PrivateKeyPwd string
	Port          int
	Timeout       int
}

// 带默认值的ssh client
func DefaultSSHClient() SSHClient {
	return SSHClient{
		Ip:            "",
		User:          "",
		Password:      "",
		PrivateKey:    "",
		PrivateKeyPwd: "",
		Port:          22,
		Timeout:       10,
	}
}

//
// 通用的ssh连接client
//

// 设置ssh所需的config
func (client *SSHClient) ParseConfig() (ssh.ClientConfig, error) {
	var (
		authMethods []ssh.AuthMethod
		signer      ssh.Signer
	)

	// 使用密码验证
	if client.Password != "" {
		authMethods = append(authMethods, ssh.Password(client.Password))
	}

	// 使用秘钥验证
	if client.PrivateKey != "" {
		fp, err := os.Open(client.PrivateKey)
		if err != nil {
			return ssh.ClientConfig{}, err
		}
		defer fp.Close()
		buf, err := ioutil.ReadAll(fp)
		if err != nil {
			return ssh.ClientConfig{}, err
		}
		if client.PrivateKeyPwd == "" {
			// 不带密码的秘钥
			signer, err = ssh.ParsePrivateKey(buf)
		} else {
			// 带密码的秘钥
			signer, err = ssh.ParsePrivateKeyWithPassphrase(buf, []byte(client.PrivateKeyPwd))
		}
		if err != nil {
			return ssh.ClientConfig{}, err
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	return ssh.ClientConfig{
		User:    client.User,
		Auth:    authMethods,
		Timeout: time.Duration(client.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}, nil
}

// 建立连接
func (s *SSHClient) Connect(config ssh.ClientConfig) (*ssh.Session, error) {

	// 建立连接
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port), &config)
	if err != nil {
		return &ssh.Session{}, err
	}

	// 建立session
	session, err := conn.NewSession()
	if err != nil {
		return &ssh.Session{}, err
	}

	return session, nil
}

// 执行命令
func (s *SSHClient) RunCmd(command string) Result {

	// 生成ssh config
	config, err := s.ParseConfig()
	if err != nil {
		return Result{
			Rst:  false,
			Msg:  fmt.Sprintf("%s %s", s.Ip, err.Error()),
			Code: 1,
		}
	}

	session, err := s.Connect(config)
	if err != nil {
		return Result{
			Rst:  false,
			Msg:  fmt.Sprintf("%s %s", s.Ip, err.Error()),
			Code: 1,
		}
	}

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)

	return Result{
		Rst:  true,
		Msg:  fmt.Sprintf("%s %s", s.Ip, stdoutBuf.String()),
		Code: 0,
	}
}

//
// 针对PMS编写合适的函数
//

// 设置ssh所需的config
func CreateSSHConfig(user, password, privateKey, keyPwd string, timeOut int) (ssh.ClientConfig, error) {
	var (
		authMethods []ssh.AuthMethod
		signer      ssh.Signer
	)

	// 使用密码验证
	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	// 使用秘钥验证
	if privateKey != "" {
		fp, err := os.Open(privateKey)
		if err != nil {
			return ssh.ClientConfig{}, err
		}
		defer fp.Close()
		buf, err := ioutil.ReadAll(fp)
		if err != nil {
			return ssh.ClientConfig{}, err
		}
		if keyPwd == "" {
			// 不带密码的秘钥
			signer, err = ssh.ParsePrivateKey(buf)
		} else {
			// 带密码的秘钥
			signer, err = ssh.ParsePrivateKeyWithPassphrase(buf, []byte(keyPwd))
		}
		if err != nil {
			return ssh.ClientConfig{}, err
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	return ssh.ClientConfig{
		User:    user,
		Auth:    authMethods,
		Timeout: time.Duration(timeOut) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}, nil
}

// 测试是否能连接
func DialSuccess(ip string, port int, config ssh.ClientConfig) Result {

	// 测试联通性
	if _, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), &config); err != nil {
		return Result{
			Rst:  false,
			Msg:  fmt.Sprintf("error to dialing ssh: %v", err),
			Code: 30001,
			Ip:   ip,
		}
	}

	return Result{
		Rst:  true,
		Msg:  "succ to dialing ssh",
		Code: 0,
		Ip:   ip,
	}
}

// 执行命令
func RunCmd(ip string, port int, command string, config ssh.ClientConfig) Result {

	// 建立连接
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), &config)
	if err != nil {
		return Result{
			Rst:  false,
			Msg:  fmt.Sprintf("error to dialing ssh: %v", err),
			Code: 30001,
			Ip:   ip,
		}
	}

	// 建立session
	session, err := conn.NewSession()
	if err != nil {
		return Result{
			Rst:  false,
			Msg:  fmt.Sprintf("error to create session: %v", err),
			Code: 30002,
			Ip:   ip,
		}
	}

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)

	return Result{
		Rst:  true,
		Msg:  stdoutBuf.String(),
		Code: 0,
		Ip:   ip,
	}
}
