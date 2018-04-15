package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"io"
)

var (
	ErrAppIDNotMatch       = errors.New("app id not match")
	ErrInvalidBlockSize    = errors.New("invalid block size")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

func Salt() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return Sha1(token.String()), nil
}

func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1Bytes(str string) []byte {
	h := sha1.New()
	io.WriteString(h, str)
	return h.Sum(nil)
}

func Sha256(str string) string {
	h := sha256.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha256Bytes(str string) []byte {
	h := sha256.New()
	io.WriteString(h, str)
	return h.Sum(nil)
}

func Hmac256(message string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// fixLengthKeyHash will transform the key ([]byte) into the first its SHA256
// checksum. This forces users to use this software to decrypt, but will
// increase security allowing them to any key length.
func fixLengthKeyHash(bytes []byte) []byte {
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}

// AESCBCEncrypt encrypts string to base64 crypto using AES
func AESCBCEncrypt(key []byte, text string) (string, error) {
	plaintext := []byte(text)

	key = fixLengthKeyHash(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// AESCBCDecrypt from base64 to decrypted string
func AESCBCDecrypt(key []byte, cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	key = fixLengthKeyHash(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)

	return string(ciphertext[:]), nil
}

// AESEncrypt encrypts string to base64 crypto using AES
func AESEncrypt(key []byte, text string) (string, error) {
	plaintext := []byte(text)

	key = fixLengthKeyHash(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt from base64 to decrypted string
func AESDecrypt(key []byte, cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	key = fixLengthKeyHash(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext[:]), nil
}

// AESEncrypt encrypts string to base64 crypto using AES
func AESEncryptBytes(key []byte, plaintext []byte) (string, error) {

	key = fixLengthKeyHash(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt from base64 to decrypted string
func AESDecryptBytes(key []byte, cryptoText string) ([]byte, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	key = fixLengthKeyHash(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext[:], nil
}

func AddressEncrypt(addr string, key string) (salt string, encrypted string, err error) {
	salt, _ = Salt()
	secret := Sha1Bytes(fmt.Sprintf("%s%s%s", key, salt, key))
	encrypted, err = AESEncrypt(secret, addr)
	return
}

func AddressDecrypt(encrypted string, salt string, key string) (addr string, err error) {
	secret := Sha1Bytes(fmt.Sprintf("%s%s%s", key, salt, key))
	addr, err = AESDecrypt(secret, encrypted)
	return
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize //需要padding的数目
	//只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	//最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) //生成填充的文本
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS7UnPadding(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
