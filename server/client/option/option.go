// Copyright 2021 Tencent Inc. All rights reserved.

package option

import (
	"cooller/server/client"
	auth2 "cooller/server/client/auth"
	"cooller/server/client/auth/signers"
	validators2 "cooller/server/client/auth/validators"
	"cooller/server/client/auth/verifiers"
	cipher2 "cooller/server/client/cipher"
	"cooller/server/client/cipher/ciphers"
	"crypto/rsa"
	"crypto/x509"
	"net/http"
)

// region SignerOption

// withSignerOption 为 Client 设置 Signer
type withSignerOption struct {
	Signer auth2.Signer
}

// Apply 将配置添加到 core.DialSettings 中
func (w withSignerOption) Apply(o *client.DialSettings) error {
	o.Signer = w.Signer
	return nil
}

// WithSigner 返回一个指定signer的ClientOption
func WithSigner(signer auth2.Signer) client.ClientOption {
	return withSignerOption{Signer: signer}
}

// WithMerchantCredential 通过商户号、商户证书序列号、商户私钥构建一对 Credential/Signer，用于生成请求头中的 Authorization 信息
func WithMerchantCredential(mchID, certificateSerialNo string, privateKey *rsa.PrivateKey) client.ClientOption {
	signer := &signers.SHA256WithRSASigner{
		MchID:               mchID,
		PrivateKey:          privateKey,
		CertificateSerialNo: certificateSerialNo,
	}
	return WithSigner(signer)
}

// endregion

// region ValidatorOption

// withValidatorOption 为 Client 设置 Validator
type withValidatorOption struct {
	Validator auth2.Validator
}

// Apply 将配置添加到 core.DialSettings 中
func (w withValidatorOption) Apply(o *client.DialSettings) error {
	o.Validator = w.Validator
	return nil
}

// WithVerifier 返回一个指定verifier的ClientOption，用于校验http response header
func WithVerifier(verifier auth2.Verifier) client.ClientOption {
	validator := validators2.NewWechatPayResponseValidator(verifier)
	return withValidatorOption{Validator: validator}
}

// WithWechatPayCertificate 设置微信支付平台证书信息，返回一个指定validator的ClientOption，用于校验http response header
func WithWechatPayCertificate(certificateList []*x509.Certificate) client.ClientOption {
	verifier := verifiers.NewSHA256WithRSAVerifier(client.NewCertificateMapWithList(certificateList))
	return WithVerifier(verifier)
}

// WithoutValidator 返回一个指定validator的ClientOption，不进行验签 用于下载证书和下载账单等不需要进行验签的接口
func WithoutValidator() client.ClientOption {
	return withValidatorOption{Validator: &validators2.NullValidator{}}
}

// endregion

// region HTTPClientOption

// withHTTPClientOption 为 Client 设置 HTTPClient
type withHTTPClientOption struct {
	Client *http.Client
}

// Apply 将配置添加到 core.DialSettings 中
func (w withHTTPClientOption) Apply(o *client.DialSettings) error {
	o.HTTPClient = w.Client
	return nil
}

// WithHTTPClient 返回一个指定网络通信为HttpClient的ClientOption，指定后使用用户自动创建的的http.client，如果用户不创建，则帮助用户
// 创建一个默认的http.client
func WithHTTPClient(client *http.Client) client.ClientOption {
	return withHTTPClientOption{Client: client}
}

// endregion

// region CipherOption

// withCipherOption 为 Client 设置 Cipher
type withCipherOption struct {
	Cipher cipher2.Cipher
}

// Apply 将配置添加到 core.DialSettings 中
func (w withCipherOption) Apply(o *client.DialSettings) error {
	o.Cipher = w.Cipher
	return nil
}

// WithWechatPayCipher 返回一个为 Client 设置 WechatPayCipher 的 core.ClientOption
func WithWechatPayCipher(encryptor cipher2.Encryptor, decryptor cipher2.Decryptor) client.ClientOption {
	return withCipherOption{Cipher: ciphers.NewWechatPayCipher(encryptor, decryptor)}
}

// endregion
