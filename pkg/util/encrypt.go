package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"time"
)

const (
	NUmStr    = "0123456789"
	CharStr   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	SpecStr   = "+=-@#~,.[]()!%^*$"
	CommonKey = "95PmsApi20180704"
)

// 加密api用户密码
func EncryptPass(username, password, salt string) string {
	// 计算密码MD5
	c_user := md5.New()
	io.WriteString(c_user, username)
	spw_user := fmt.Sprintf("%x", c_user.Sum(nil))

	c_passwd := md5.New()
	io.WriteString(c_passwd, password)
	spw_passwd := fmt.Sprintf("%x", c_passwd.Sum(nil))

	// 拼接密码MD5
	buf := bytes.NewBufferString("")

	// 拼接密码
	io.WriteString(buf, spw_passwd)
	io.WriteString(buf, salt)
	io.WriteString(buf, spw_user)

	// 拼接密码计算MD5
	t := md5.New()
	io.WriteString(t, buf.String())

	return fmt.Sprintf("%x", t.Sum(nil))
}

// 生成系统用户密码
func GeneratePasswd() string {
	// 8位密码
	length := 8
	var passwd []byte = make([]byte, length, length)
	sourceStr := fmt.Sprintf("%s%s%s", NUmStr, CharStr, SpecStr)
	//随机种子
	rand.Seed(time.Now().UnixNano())
	// 遍历，生成一个随机index索引
	for i := 0; i < length; i++ {
		index := rand.Intn(len(sourceStr))
		passwd[i] = sourceStr[index]
	}

	return string(passwd)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 加密系统用户密码
func EncryptUserPwd(passwd []byte) (string, error) {

	block, err := aes.NewCipher([]byte(CommonKey))
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	passwd = PKCS7Padding(passwd, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(CommonKey)[:blockSize])
	crypted := make([]byte, len(passwd))
	blockMode.CryptBlocks(crypted, passwd)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

// 解密系统用户密码
func DecryptUserPwd(crypted string) (string, error) {

	ncryted, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(CommonKey))
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, []byte(CommonKey)[:blockSize])
	origData := make([]byte, len(ncryted))
	blockMode.CryptBlocks(origData, ncryted)
	origData = PKCS7UnPadding(origData)

	return string(origData), nil
}
