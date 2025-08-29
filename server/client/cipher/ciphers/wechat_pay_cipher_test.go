// Copyright 2021 Tencent Inc. All rights reserved.

package ciphers

import (
	"context"
	"cooller/server/client"
	"cooller/server/client/cipher/decryptors"
	encryptors2 "cooller/server/client/cipher/encryptors"
	"cooller/server/utils"
	"github.com/agiledragon/gomonkey"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Student struct {
	Name      string `encryption:"EM_APIV3"`
	Age       int
	Addresses []Address
	Parents   *[]Parent
	// unexported secret
	secret string `encryption:"EM_APIV3"`
	IDs    []int
}

type Address struct {
	// No Tag
	Country *string
	// Not EM_APIV3 encryption Tag
	Province *string `encryption:"EM_APIV2"`
	// EM_APIV3 encryption Tag
	City   **string `encryption:"EM_APIV3"`
	Street *string  `encryption:"EM_APIV3"`
}

type Parent struct {
	Name        string  `encryption:"EM_APIV3"`
	PhoneNumber *string `encryption:"EM_APIV3"`
}

func TestContextKey_String(t *testing.T) {
	assert.Equal(t, "WPCipherContext(EncryptSerial)", contextKeyEncryptSerial.String())
}

func TestWechatPayCipher_Encrypt_Decrypt(t *testing.T) {
	cityCD := utils.String("成都")
	cityLA := utils.String("LA")

	s := Student{
		Name: "小可",
		Age:  8,
		Addresses: []Address{
			{
				Country:  utils.String("中国"),
				Province: utils.String("四川"),
				City:     &cityCD,
				Street:   utils.String("春熙路"),
			},
			{
				Country:  utils.String("USA"),
				Province: utils.String("California"),
				City:     &cityLA,
				Street:   utils.String("Nowhere"),
			},
		},
		Parents: &[]Parent{
			{
				Name:        "爸",
				PhoneNumber: utils.String("13000000000"),
			},
			{
				Name:        "妈",
				PhoneNumber: nil,
			},
		},
		secret: "this is secret",
		IDs: []int{
			12345,
			54321,
		},
	}

	c := WechatPayCipher{
		encryptor: &encryptors2.MockEncryptor{
			Serial: "Mock Serial",
		},
		decryptor: &decryptors.MockDecryptor{},
	}

	serial, err := c.Encrypt(context.Background(), &s)
	assert.Equal(t, "Mock Serial", serial)
	require.NoError(t, err)
	assert.Equal(t, "Encrypted小可", s.Name)
	assert.Equal(t, 8, s.Age)
	assert.Equal(t, "中国", *(s.Addresses[0].Country))
	assert.Equal(t, "四川", *(s.Addresses[0].Province))
	assert.Equal(t, "Encrypted成都", **(s.Addresses[0].City))
	assert.Equal(t, "Encrypted春熙路", *(s.Addresses[0].Street))
	assert.Equal(t, "USA", *(s.Addresses[1].Country))
	assert.Equal(t, "California", *(s.Addresses[1].Province))
	assert.Equal(t, "EncryptedLA", **(s.Addresses[1].City))
	assert.Equal(t, "EncryptedNowhere", *(s.Addresses[1].Street))
	assert.Equal(t, "Encrypted爸", (*s.Parents)[0].Name)
	assert.Equal(t, "Encrypted13000000000", *((*s.Parents)[0].PhoneNumber))
	assert.Equal(t, "Encrypted妈", (*s.Parents)[1].Name)
	assert.Equal(t, (*string)(nil), (*s.Parents)[1].PhoneNumber)
	assert.Equal(t, "this is secret", s.secret) // unexported fields will be skipped
	assert.Equal(t, 12345, s.IDs[0])
	assert.Equal(t, 54321, s.IDs[1])

	err = c.Decrypt(context.Background(), &s)
	require.NoError(t, err)
	assert.Equal(t, "小可", s.Name)
	assert.Equal(t, 8, s.Age)
	assert.Equal(t, "中国", *(s.Addresses[0].Country))
	assert.Equal(t, "四川", *(s.Addresses[0].Province))
	assert.Equal(t, "成都", **(s.Addresses[0].City))
	assert.Equal(t, "春熙路", *(s.Addresses[0].Street))
	assert.Equal(t, "USA", *(s.Addresses[1].Country))
	assert.Equal(t, "California", *(s.Addresses[1].Province))
	assert.Equal(t, "LA", **(s.Addresses[1].City))
	assert.Equal(t, "Nowhere", *(s.Addresses[1].Street))
	assert.Equal(t, "爸", (*s.Parents)[0].Name)
	assert.Equal(t, "13000000000", *((*s.Parents)[0].PhoneNumber))
	assert.Equal(t, "妈", (*s.Parents)[1].Name)
	assert.Equal(t, (*string)(nil), (*s.Parents)[1].PhoneNumber)
	assert.Equal(t, "this is secret", s.secret) // unexported fields will be skipped
	assert.Equal(t, 12345, s.IDs[0])
	assert.Equal(t, 54321, s.IDs[1])
}

func TestWechatPayCipher_Encrypt_DecryptWithValue(t *testing.T) {
	cityCD := utils.String("成都")
	cityLA := utils.String("LA")

	s := Student{
		Name: "小可",
		Age:  8,
		Addresses: []Address{
			{
				Country:  utils.String("中国"),
				Province: utils.String("四川"),
				City:     &cityCD,
				Street:   utils.String("春熙路"),
			},
			{
				Country:  utils.String("USA"),
				Province: utils.String("California"),
				City:     &cityLA,
				Street:   utils.String("Nowhere"),
			},
		},
		Parents: &[]Parent{
			{
				Name:        "爸",
				PhoneNumber: utils.String("13000000000"),
			},
			{
				Name:        "妈",
				PhoneNumber: nil,
			},
		},
	}

	c := NewWechatPayCipher(
		&encryptors2.MockEncryptor{
			Serial: "Mock Serial",
		},
		&decryptors.MockDecryptor{},
	)

	serial, err := c.Encrypt(context.Background(), reflect.ValueOf(&s))
	assert.Equal(t, "Mock Serial", serial)
	require.NoError(t, err)
	assert.Equal(t, "Encrypted小可", s.Name)
	assert.Equal(t, 8, s.Age)
	assert.Equal(t, "中国", *(s.Addresses[0].Country))
	assert.Equal(t, "四川", *(s.Addresses[0].Province))
	assert.Equal(t, "Encrypted成都", **(s.Addresses[0].City))
	assert.Equal(t, "Encrypted春熙路", *(s.Addresses[0].Street))
	assert.Equal(t, "USA", *(s.Addresses[1].Country))
	assert.Equal(t, "California", *(s.Addresses[1].Province))
	assert.Equal(t, "EncryptedLA", **(s.Addresses[1].City))
	assert.Equal(t, "EncryptedNowhere", *(s.Addresses[1].Street))
	assert.Equal(t, "Encrypted爸", (*s.Parents)[0].Name)
	assert.Equal(t, "Encrypted13000000000", *((*s.Parents)[0].PhoneNumber))
	assert.Equal(t, "Encrypted妈", (*s.Parents)[1].Name)
	assert.Equal(t, (*string)(nil), (*s.Parents)[1].PhoneNumber)

	err = c.Decrypt(context.Background(), reflect.ValueOf(&s))
	require.NoError(t, err)
	assert.Equal(t, "小可", s.Name)
	assert.Equal(t, 8, s.Age)
	assert.Equal(t, "中国", *(s.Addresses[0].Country))
	assert.Equal(t, "四川", *(s.Addresses[0].Province))
	assert.Equal(t, "成都", **(s.Addresses[0].City))
	assert.Equal(t, "春熙路", *(s.Addresses[0].Street))
	assert.Equal(t, "USA", *(s.Addresses[1].Country))
	assert.Equal(t, "California", *(s.Addresses[1].Province))
	assert.Equal(t, "LA", **(s.Addresses[1].City))
	assert.Equal(t, "Nowhere", *(s.Addresses[1].Street))
	assert.Equal(t, "爸", (*s.Parents)[0].Name)
	assert.Equal(t, "13000000000", *((*s.Parents)[0].PhoneNumber))
	assert.Equal(t, "妈", (*s.Parents)[1].Name)
	assert.Equal(t, (*string)(nil), (*s.Parents)[1].PhoneNumber)
}

func TestWechatPayCipher_CipherNil(t *testing.T) {
	c := WechatPayCipher{
		encryptor: &encryptors2.MockEncryptor{
			Serial: "Mock Serial",
		},
		decryptor: &decryptors.MockDecryptor{},
	}

	var s *Student

	_, err := c.Encrypt(context.Background(), s)
	require.NoError(t, err)

	err = c.Decrypt(context.Background(), &s)
	require.NoError(t, err)
}

func TestWechatPayCipher_CipherNonStruct(t *testing.T) {
	c := WechatPayCipher{
		encryptor: &encryptors2.MockEncryptor{
			Serial: "Mock Serial",
		},
		decryptor: &decryptors.MockDecryptor{},
	}

	_, err := c.Encrypt(context.Background(), utils.String("123"))
	require.Error(t, err)
	assert.Equal(t, "encrypt struct failed: only struct can be ciphered", err.Error())

	err = c.Decrypt(context.Background(), utils.Int64(123))
	require.Error(t, err)
	assert.Equal(t, "decrypt struct failed: only struct can be ciphered", err.Error())
}

func TestWechatPayCipher_CipherValue(t *testing.T) {
	s := Student{
		Name: "小可",
		Age:  8,
	}

	c := WechatPayCipher{
		encryptor: &encryptors2.MockEncryptor{
			Serial: "Mock Serial",
		},
		decryptor: &decryptors.MockDecryptor{},
	}

	_, err := c.Encrypt(context.Background(), s)
	require.Error(t, err)
	assert.Equal(t, "encrypt struct failed: in-place cipher requires settable input, ptr for example", err.Error())

	err = c.Decrypt(context.Background(), s)
	require.Error(t, err)
	assert.Equal(t, "decrypt struct failed: in-place cipher requires settable input, ptr for example", err.Error())
}

func TestWechatPayCipher_EncryptWithoutCertificate(t *testing.T) {
	s := Student{Name: "小可"}

	// 这是一个 SelectCertificate 会失败的 Encryptor
	invalidEncryptor := encryptors2.NewWechatPayEncryptor(client.NewCertificateMap(nil))

	c := WechatPayCipher{
		encryptor: invalidEncryptor,
		decryptor: &decryptors.MockDecryptor{},
	}

	_, err := c.Encrypt(context.Background(), s)
	assert.Error(t, err)
}

func TestWechatPayCipher_EncryptWithoutSerial(t *testing.T) {
	patch := gomonkey.ApplyFunc(getEncryptSerial, func(ctx context.Context) (string, bool) {
		return "", false
	})
	defer patch.Reset()
	s := Student{
		Name: "小可",
		Age:  8,
	}

	c := WechatPayCipher{
		encryptor: &encryptors2.MockEncryptor{
			Serial: "Mock Serial",
		},
		decryptor: &decryptors.MockDecryptor{},
	}

	_, err := c.Encrypt(context.Background(), &s)
	assert.Error(t, err)
}

func TestWechatPayCipher_DecryptWrongData(t *testing.T) {
	s := Student{
		Name: "NotEncrypted小可",
		Age:  8,
	}

	c := WechatPayCipher{
		encryptor: &encryptors2.MockEncryptor{
			Serial: "Mock Serial",
		},
		decryptor: &decryptors.MockDecryptor{},
	}

	err := c.Decrypt(context.Background(), &s)
	assert.Error(t, err)

	s = Student{
		Name: "Encrypted小可",
		Addresses: []Address{
			{
				Country:  utils.String("中国"),
				Province: utils.String("四川"),
				Street:   utils.String("UnEncrypted春熙路"),
			},
			{
				Country:  utils.String("USA"),
				Province: utils.String("California"),
				Street:   utils.String("EncryptedNowhere"),
			},
		},
	}

	err = c.Decrypt(context.Background(), &s)
	assert.Error(t, err)
}

func TestWechatPayCipher_cipherWithWrongType(t *testing.T) {
	s := Student{
		Name: "Encrypted小可",
		Age:  8,
	}

	c := WechatPayCipher{
		encryptor: &encryptors2.MockEncryptor{
			Serial: "Mock Serial",
		},
		decryptor: &decryptors.MockDecryptor{},
	}

	err := c.cipher(context.Background(), cipherType("invalid"), reflect.ValueOf(&s))
	assert.Error(t, err)
}
