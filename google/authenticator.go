package google

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

type Authenticator struct {
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

func (ga *Authenticator) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (ga *Authenticator) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (ga *Authenticator) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (ga *Authenticator) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (ga *Authenticator) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (ga *Authenticator) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (ga *Authenticator) oneTimePassword(key []byte, data []byte) uint32 {
	hash := ga.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := ga.toUint32(hashParts)
	return number % 1000000
}

// GetSecret 获取秘钥
func (ga *Authenticator) GetSecret() string {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, ga.un())
	return strings.ToUpper(ga.base32encode(ga.hmacSha1(buf.Bytes(), nil)))
}

// GetCode 获取动态码
func (ga *Authenticator) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := ga.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := ga.oneTimePassword(secretKey, ga.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}

// GetQrcode 获取动态码二维码内容
func (ga *Authenticator) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

// GetQrcodeUrl 获取动态码二维码图片地址,这里是第三方二维码api
func (ga *Authenticator) GetQrcodeUrl(user, secret string, Status ...int) string {
	size := 200
	if len(Status) > 0 {
		size = Status[0]
	}
	_qrcode := ga.GetQrcode(user, secret)

	qrByte, _ := qrcode.Encode(_qrcode, qrcode.Medium, size)
	res := base64.StdEncoding.EncodeToString(qrByte)
	return fmt.Sprintf("data:image/jpg;base64,%s", res)
}

// VerifyCode 验证动态码
func (ga *Authenticator) VerifyCode(secret, code string) (bool, error) {
	_code, err := ga.GetCode(secret)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}
