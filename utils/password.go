package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"github.com/astaxie/beego/logs"
	"io"
	mt "math/rand"
	"strconv"
	"strings"
)

const (
	saltSize            = 16
	delmiter            = "$"
	stretching_password = 500
	salt_local_secret   = "ahefew*&TGEsfdbi*^WB"
	PrivateKey          = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDetmYzEUOTog5zPtp29OByTcDQPgusyEx9DzKYJZNTzRiX/22u
t0OIgtfVE80vz++JwXUljzafk1JD6QmwE+lp0erPow6PqYCP1vX29fbhUW/9nRe+
vCnjLU7JoxrIsBjSFaN4YoWtKFZPRts5dQZw8UFiw5KCMp0WD74AdhvlwwIDAQAB
AoGAVGAsNfq7bGpAKT9NyzWY9xUoEH0BNVOpTtP8KhJKT7xrLeLSrhe2WTihBpP6
77tKmBkYBcPNQQWybBIU3oWcrwpwLJJbVDmmpYQ6VRAjDo2EXrNbJWUbpBh63pTs
h5JPPhDovMqAcklrNuYYPbUEqaTQAED1ysf2RGLfey3ZupkCQQD7rpF5v3S11hR4
cU644edcbUxo5JGF+ULF7E/CwSuHQFwrfrU/HYXodeNGl8tJE+I9ZTxXY9vlERmr
0Fh84uMXAkEA4oiXluFFhjCr8YGvS+Wc8J54vw2mHDpJEozPefnIZKtL8FzlVfo+
zNZDQs5bPSblqr0ipjrhbQ4nFbQqBEluNQJBAJ7CE0n1Fy3MiMUg1EOTXFnVKCnS
ZGlaPmCTHA0BxO9gDcPx/Wp+uQVVt7PD9Jt4S3Hm9hU6DG+GRec3WVoN1KkCQQCc
l91qIAkGTOjfBk2eAnhtYK6JKy8zfhr7JrlZURB0fnD9E8o4l8cHo+lU6f7qE9RZ
JWspS7R+xXTBLQyKcBQtAkBvTCisqbpmtiXET2yWJMFAaDMIq+LUxm0y/9AS3OZk
yQpEmeFECWoCGyh7R/X3NL4GHjyr2X7XzWbp/RPxlgor
-----END RSA PRIVATE KEY-----`

	PublickKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDetmYzEUOTog5zPtp29OByTcDQ
PgusyEx9DzKYJZNTzRiX/22ut0OIgtfVE80vz++JwXUljzafk1JD6QmwE+lp0erP
ow6PqYCP1vX29fbhUW/9nRe+vCnjLU7JoxrIsBjSFaN4YoWtKFZPRts5dQZw8UFi
w5KCMp0WD74AdhvlwwIDAQAB
-----END PUBLIC KEY-----`
)

//加密密码
func PasswordHash(pass string) (string, error) {

	salt_secret, err := salt_secret()
	if err != nil {
		return "", err
	}

	salt, err := salt(salt_local_secret + salt_secret)
	if err != nil {
		return "", err
	}

	interation := randInt(1, 20)

	hash, err := hash(pass, salt_secret, salt, int64(interation))
	if err != nil {
		return "", err
	}
	interation_string := strconv.Itoa(interation)
	password := salt_secret + delmiter + interation_string + delmiter + hash + delmiter + salt

	return password, nil

}

//校验密码是否有效
func PasswordVerify(hashing string, pass string) (bool, error) {

	defer func() {
		if err := recover(); err != nil {
			logs.Info("PasswordVerify %s", err)
		}
	}()

	data := trim_salt_hash(hashing)

	interation, _ := strconv.ParseInt(data["interation_string"], 10, 64)

	has, err := hash(pass, data["salt_secret"], data["salt"], int64(interation))
	if err != nil {
		return false, err
	}

	if (data["salt_secret"] + delmiter + data["interation_string"] + delmiter + has + delmiter + data["salt"]) == hashing {
		return true, nil
	} else {
		return false, nil
	}

}

func hash(pass string, salt_secret string, salt string, interation int64) (string, error) {
	var pass_salt string = salt_secret + pass + salt + salt_secret + pass + salt + pass + pass + salt
	var i int

	hash_pass := salt_local_secret
	hash_start := sha512.New()
	hash_center := sha256.New()
	hash_output := sha256.New224()

	i = 0
	for i <= stretching_password {
		i = i + 1
		hash_start.Write([]byte(pass_salt + hash_pass))
		hash_pass = hex.EncodeToString(hash_start.Sum(nil))
	}

	i = 0
	for int64(i) <= interation {
		i = i + 1
		hash_pass = hash_pass + hash_pass
	}

	i = 0
	for i <= stretching_password {
		i = i + 1
		hash_center.Write([]byte(hash_pass + salt_secret))
		hash_pass = hex.EncodeToString(hash_center.Sum(nil))
	}
	hash_output.Write([]byte(hash_pass + salt_local_secret))
	hash_pass = hex.EncodeToString(hash_output.Sum(nil))

	return hash_pass, nil
}

func trim_salt_hash(hash string) map[string]string {
	str := strings.Split(hash, delmiter)

	return map[string]string{
		"salt_secret":       str[0],
		"interation_string": str[1],
		"hash":              str[2],
		"salt":              str[3],
	}
}
func salt(secret string) (string, error) {

	buf := make([]byte, saltSize, saltSize+md5.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write(buf)
	hash.Write([]byte(secret))
	return hex.EncodeToString(hash.Sum(buf)), nil
}

func salt_secret() (string, error) {
	rb := make([]byte, randInt(10, 100))
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

func randInt(min int, max int) int {
	return min + mt.Intn(max-min)
}

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(PublickKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(PrivateKey))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
