package serial

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash/crc32"
	"os"
)

const gEncodeLen = 20
const gDecodeLen = 32

// 从文件中读取取公钥
func GetPublicKey(pubFile string) (pubkey rsa.PublicKey, err error) {
	stat, err := os.Stat(pubFile)
	if err != nil {
		return pubkey, err
	}

	fp, err := os.Open(pubFile)
	defer fp.Close()

	key := make([]byte, stat.Size())
	_, err = fp.Read(key)
	if err != nil {
		return pubkey, err
	}

	jsonPubKey, err := base64.StdEncoding.DecodeString(string(key))
	err = json.Unmarshal(jsonPubKey, &pubkey)

	return pubkey, err
}

// 使用公钥加密数据
func RsaEncrypt(pubKey *rsa.PublicKey, origData []byte) (enData []byte, err error) {
	msgLen := len(origData)
	var en []byte
	for i := 0; i < msgLen; i += gEncodeLen {
		if i+gEncodeLen > msgLen {
			en, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, origData[i:msgLen])
		} else {
			en, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, origData[i:i+gEncodeLen])
		}
		if err != nil {
			return enData, err
		}

		enData = append(enData, en...)
	}

	return enData, err
	//return rsa.EncryptPKCS1v15(rand.Reader, pubKey, origData)
}

// 从文件中读取私钥
func GetPrivateKey(priFile string) (privkey rsa.PrivateKey, err error) {
	stat, err := os.Stat(priFile)
	if err != nil {
		return privkey, err
	}

	fp, err := os.Open(priFile)
	defer fp.Close()

	key := make([]byte, stat.Size())
	_, err = fp.Read(key)
	if err != nil {
		return privkey, err
	}
	jsonPriKey, err := base64.StdEncoding.DecodeString(string(key))
	err = json.Unmarshal(jsonPriKey, &privkey)

	return privkey, err
}

// 使用私钥解密数据
func RsaDecrypt(privateKey *rsa.PrivateKey, encodeData []byte) (deData []byte, err error) {
	msgLen := len(encodeData)
	if (msgLen % gDecodeLen) != 0 {
		return deData, errors.New("加密数据长度不合法")
	}
	for i := 0; i < msgLen; i += gDecodeLen {
		de, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encodeData[i:i+gDecodeLen])
		if err != nil {
			return deData, err
		}

		deData = append(deData, de...)
	}

	// 去掉末尾对齐的\0x00
	deLen := len(deData)
	var i int
	for i = deLen - 1; i > 0; i-- {
		if deData[i] != 0 {
			break
		}
	}
	return deData[:i+1], err
}

// 使用私钥签名
func RsaSign(privateKey *rsa.PrivateKey, sigData []byte) (sig []byte, err error) {
	var digest []byte

	digest = GetShA1(sigData)
	sig, err = rsa.SignPSS(rand.Reader, privateKey, crypto.SHA1, digest, nil)

	return sig, err
}

// 公钥验证签名
func RsaSignVerify(pubKey *rsa.PublicKey, sigData []byte, sig []byte) (err error) {
	var digest []byte

	digest = GetShA1(sigData)
	err = rsa.VerifyPSS(pubKey, crypto.SHA1, digest, sig, nil)

	return err
}

// 获取SHA1
func GetShA1(data []byte) (digest []byte) {
	h := sha1.New()
	h.Write(data)
	digest = h.Sum(nil)
	return digest
}

// 获取CRC32
func GetCrc32(data []byte) (crc uint32) {
	h := crc32.NewIEEE()
	h.Write(data)
	return h.Sum32()
}
