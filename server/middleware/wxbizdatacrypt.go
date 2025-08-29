package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
)

// DecryptWXOpenData 微信小程序解密算法 AES-128-CBC
func DecryptWXOpenData(sessionKey, encryptData, iv string) (map[string]interface{}, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return nil, err
	}
	sessionKeyBytes, errKey := base64.StdEncoding.DecodeString(sessionKey)
	if errKey != nil {
		return nil, errKey
	}
	ivBytes, errIv := base64.StdEncoding.DecodeString(iv)
	if errIv != nil {
		return nil, errIv
	}
	dataBytes, errData := aesDecrypt(decodeBytes, sessionKeyBytes, ivBytes)
	if errData != nil {
		return nil, errData
	}

	var result map[string]interface{}
	errResult := json.Unmarshal(dataBytes, &result)

	return result, errResult
}

// AES 解密
func aesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	// 原始数据
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	// 去除填充  --- 数据尾端有'/x0e'占位符,去除它
	length := len(origData)
	unp := int(origData[length-1])
	return origData[:(length - unp)], nil
}
