package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strconv"
	"strings"
	"bytes"
	"crypto/des"
	"crypto/cipher"
	"crypto/x509"
	"crypto/rsa"
	"crypto"
	"fmt"
)

const (
	//BASE64字符表,不要有重复
	base64Table        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

//默认把时间戳变为二进制字符串
func GetUserValCodeDefault(userID string, time int64) string {
	return GetUserValCode(userID, time, 2)
}

func GetUserValCode(userID string, time int64, num int) string {
	timeStr := strconv.FormatInt(time, num)
	return GetMd5String(timeStr)
}

/**
 * 对一个字符串进行MD5加密,不可解密
 */
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s)) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))
}

/*获取 SHA1 字符串*/
func GetSHA1String(s string) string {
	t := sha1.New()
	t.Write([]byte(s))
	return hex.EncodeToString(t.Sum(nil))
}

/**
 * 获取一个Guid值
 */
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

var coder = base64.NewEncoding(base64Table)

/**
 * base64加密
 */
func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

/**
 * base64解密
 */
func Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

//byte转16进制字符串
func ByteToHex(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {

		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

//16进制字符串转[]byte
func HexToBye(hex string) []byte {
	length := len(hex) / 2
	slice := make([]byte, length)
	rs := []rune(hex)

	for i := 0; i < length; i++ {
		s := string(rs[i*2 : i*2+2])
		value, _ := strconv.ParseInt(s, 16, 10)
		slice[i] = byte(value & 0xFF)
	}
	return slice
}

//判断字符串是否为空
func IsNull(str string) bool {
	if len(strings.TrimSpace(str)) != 0 {
		return false
	}
	return true
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

//DES
func DESEncode(stringStr,keyStr string)(string,error){
	key := []byte(strings.ToUpper(Substr(GetMd5String(keyStr),0,8)))
	result, err := DesEncrypt([]byte(stringStr), key)
	if err != nil {
		return "",err
	}
	hexStr:=ByteToHex(result)
	return strings.ToUpper(hexStr),nil
}

func DESDecode(hexStr,keyStr string)(string,error){
	hexByte:=HexToBye(hexStr)
	key := []byte(strings.ToUpper(Substr(GetMd5String(keyStr),0,8)))
	stringStr,err:=DesDecrypt(hexByte,key)
	if err != nil {
		return "",err
	}
	return string(stringStr),nil
}

func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}


//合并数组
func MergeArray(s ...[]interface{}) (slice []interface{}) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
		break
	default:
		s1 := s[0]
		s2 := MergeArray(s[1:]...) //...将数组元素打散
		slice = make([]interface{}, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
		break
	}

	return
}

//数组删除元素
func SliceRemove(ss []interface{}, s interface{}) []interface{} {
	for i := range ss {
		if ss[i] == s {
			ss = append(ss[:i], ss[i+1:]...)
		}
	}
	return ss
}

//rsa私钥加密
func Sha1WithRSAPKCS8Base64Sign(data string, privatekey string) (string, error) {
	key, _ := base64.StdEncoding.DecodeString(privatekey)
	privateKey, _ := x509.ParsePKCS8PrivateKey(key)
	h := sha1.New()
	h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := base64.StdEncoding.EncodeToString(signature)
	return out, nil
}

//rsa公钥解密
func Sha1WithRSAPKCS8Base64VerifySign(originData, signData, publickey string) error {
	key, _ := base64.StdEncoding.DecodeString(publickey)
	pub, _ := x509.ParsePKIXPublicKey(key)
	h := sha1.New()
	h.Write([]byte([]byte(originData)))
	hash := h.Sum(nil)
	sig, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash, sig)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return err
	}
	return nil
}