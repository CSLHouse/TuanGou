// Copyright 2021 Tencent Inc. All rights reserved.

package option

import (
	"context"
	"cooller/server/client"
	"cooller/server/client/auth/signers"
	"cooller/server/client/auth/validators"
	"cooller/server/client/auth/verifiers"
	"cooller/server/client/cipher/ciphers"
	"cooller/server/client/cipher/decryptors"
	"cooller/server/client/cipher/encryptors"
	"cooller/server/client/consts"
	downloader2 "cooller/server/client/downloader"
	"cooller/server/utils"
	"crypto/rsa"
	"crypto/x509"
	"log"
)

type withAuthCipherOption struct{ settings client.DialSettings }

// Apply 设置 core.DialSettings 的 Signer、Validator 以及 Cipher
func (w withAuthCipherOption) Apply(o *client.DialSettings) error {
	o.Signer = w.settings.Signer
	o.Validator = w.settings.Validator
	o.Cipher = w.settings.Cipher
	return nil
}

// WithWechatPayAuthCipher 一键初始化 Client，使其具备「签名/验签/敏感字段加解密」能力
func WithWechatPayAuthCipher(
	mchID string, certificateSerialNo string, privateKey *rsa.PrivateKey, certificateList []*x509.Certificate,
) client.ClientOption {
	certGetter := client.NewCertificateMapWithList(certificateList)
	return withAuthCipherOption{
		settings: client.DialSettings{
			Signer: &signers.SHA256WithRSASigner{
				MchID:               mchID,
				PrivateKey:          privateKey,
				CertificateSerialNo: certificateSerialNo,
			},
			Validator: validators.NewWechatPayResponseValidator(verifiers.NewSHA256WithRSAVerifier(certGetter)),
			Cipher: ciphers.NewWechatPayCipher(
				encryptors.NewWechatPayEncryptor(certGetter),
				decryptors.NewWechatPayDecryptor(privateKey),
			),
		},
	}
}

// WithWechatPayAutoAuthCipher 一键初始化 Client，使其具备「签名/验签/敏感字段加解密」能力。
// 同时提供证书定时更新功能（因此需要提供 mchAPIv3Key 用于证书解密），不再需要本地提供平台证书
func WithWechatPayAutoAuthCipher() client.ClientOption {
	mgr := downloader2.MgrInstance()
	if !mgr.HasDownloader(context.Background(), consts.MachID) {
		// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
		privateKey, err := utils.LoadPrivateKeyWithPath("./cert/apiclient_key.pem")
		if err != nil {
			log.Print("load merchant private key error")
		}
		// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
		err = mgr.RegisterDownloaderWithPrivateKey(
			context.Background(), privateKey, consts.MchCertificateSerialNumber, consts.MachID, consts.MchAPIv3Key,
		)
		if err != nil {
			return client.ErrorOption{Error: err}
		}
	}
	privateKey, err := utils.LoadPrivateKeyWithPath("./cert/apiclient_key.pem")
	if err != nil {
		log.Print("load merchant private key error")
	}
	return WithWechatPayAutoAuthCipherUsingDownloaderMgr(consts.MachID, consts.MchCertificateSerialNumber, privateKey, mgr)
}

// WithWechatPayAutoAuthCipherUsingDownloaderMgr 一键初始化 Client，使其具备「签名/验签/敏感字段加解密」能力。
// 需要使用者自行提供 CertificateDownloaderMgr 已实现平台证书的自动更新
//
// 【注意】本函数的能力与 WithWechatPayAutoAuthCipher 完全一致，除非有自行管理 CertificateDownloaderMgr 的需求，
// 否则推荐直接使用 WithWechatPayAutoAuthCipher
func WithWechatPayAutoAuthCipherUsingDownloaderMgr(
	mchID string, certificateSerialNo string, privateKey *rsa.PrivateKey, mgr *downloader2.CertificateDownloaderMgr,
) client.ClientOption {
	certVisitor := mgr.GetCertificateVisitor(mchID)
	return withAuthCipherOption{
		settings: client.DialSettings{
			Signer: &signers.SHA256WithRSASigner{
				MchID:               mchID,
				CertificateSerialNo: certificateSerialNo,
				PrivateKey:          privateKey,
			},
			Validator: validators.NewWechatPayResponseValidator(verifiers.NewSHA256WithRSAVerifier(certVisitor)),
			Cipher: ciphers.NewWechatPayCipher(
				encryptors.NewWechatPayEncryptor(certVisitor),
				decryptors.NewWechatPayDecryptor(privateKey),
			),
		},
	}
}
