package serial

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"../RootDir"
)

func ClientGetPubFile() string {
	subDir := "license/pub.key"
	rootDir, _ := RootDir.GetRootDir()
	if rootDir == "" {
		return "./" + subDir
	}

	return filepath.Join(rootDir, subDir)
}

func ClientGetSerialFile() (snFile string) {
	subDir := "license/licence.key"
	rootDir, _ := RootDir.GetRootDir()
	if rootDir == "" {
		return "./" + subDir
	}

	return filepath.Join(rootDir, subDir)
}

// 提取客户机注册信息
func ClientGetRegInfo() (encodeString string, err error) {
	pubFile := ClientGetPubFile()
	// 1.读取pub文件
	stat, err := os.Stat(pubFile)
	if err != nil {
		return encodeString, errors.New(fmt.Sprintf("%s 文件不存在\n", pubFile))
	}

	fp, err := os.Open(pubFile)
	defer fp.Close()

	key := make([]byte, stat.Size())
	_, err = fp.Read(key)
	if err != nil {
		return encodeString, errors.New(fmt.Sprintf("%s 读取失败\n", pubFile))
	}

	// 2.获取机器码
	var sysinfo HardWareInfo
	sysinfo, err = GetSysInfo()
	if err != nil {
		return encodeString, errors.New("获取硬件信息失败")
	}

	// 3.计算校验
	crc1 := GetCrc32(key)
	crc2 := GetCrc32([]byte(sysinfo.StaticInfo + sysinfo.CpuInfo))
	crc3 := GetCrc32([]byte(sysinfo.DiskInfo))

	// 4.拼接客户端信息
	encodeString = fmt.Sprintf("%08X-%08X-%08X", crc1, crc2, crc3)
	return encodeString, nil
}

// 从文件读取授权码
func ClientReadLicense() (serial string, err error) {
	snFile := ClientGetSerialFile()

	stat, err := os.Stat(snFile)
	if err != nil {
		return serial, err
	}

	fp, err := os.Open(snFile)
	defer fp.Close()

	bSn := make([]byte, stat.Size())
	_, err = fp.Read(bSn)
	if err != nil {
		return serial, err
	}

	return string(bSn), nil
}

// 授权码写入文件
func ClientSaveLicense(serial string) (err error) {
	snFile := ClientGetSerialFile()

	fp, err := os.Create(snFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.WriteString(serial)
	return err
}

// 验证授权文件
func ClientVerifyLicense() (err error) {
	serial, err := ClientReadLicense()
	if err != nil {
		return errors.New("错误:读取授权文件失败")
	}

	err = ClientVerifySn(serial)
	return err
}

// 验证注册码
func ClientVerifySn(serial string) (err error) {
	pubFile := ClientGetPubFile()
	// 解析注册信息
	verDaySign, err := base64.StdEncoding.DecodeString(serial)
	if err != nil {
		return errors.New("注册码格式错误")
	}

	if len(verDaySign) != 44 {
		return errors.New("注册码格式错误")
	}

	verday := string(verDaySign[0:12])
	sign := verDaySign[12:]

	// 校验有效期
	year, err := strconv.Atoi(string(verDaySign[4:8]))
	if err != nil {
		return errors.New("注册码格式错误")
	}
	mon, err := strconv.Atoi(string(verDaySign[8:10]))
	if err != nil {
		return errors.New("注册码格式错误")
	}
	day, err := strconv.Atoi(string(verDaySign[10:12]))
	if err != nil {
		return errors.New("注册码格式错误")
	}

	tm := time.Now()
	if tm.Year() > year {
		return errors.New("授权已过期")
	}

	if tm.Year() == year && int(tm.Month()) > mon {
		return errors.New("授权已过期")
	}

	if tm.Year() == year && int(tm.Month()) == mon && tm.Day() > day {
		return errors.New("授权已过期")
	}

	// 取本机硬件信息
	clientReg, err := ClientGetRegInfo()
	if err != nil {
		return err
	}

	// 校验签名
	pubKey, err := GetPublicKey(pubFile)
	if err != nil {
		return errors.New("获取密钥文件失败")
	}

	err = RsaSignVerify(&pubKey, []byte(clientReg+verday), sign)
	if err != nil {
		return errors.New("注册码不正确")
	}

	return nil
}

func ClientgetVersion() (version string, err error) {
	serial, err := ClientReadLicense()
	if err != nil {
		return version, err
	}

	verDaySign, err := base64.StdEncoding.DecodeString(serial)
	if err != nil {
		return version, err
	}

	verH, err := strconv.Atoi(string(verDaySign[0:2]))
	if err != nil {
		return version, err
	}

	verL, err := strconv.Atoi(string(verDaySign[2:4]))
	if err != nil {
		return version, err
	}

	version = fmt.Sprintf("%d.%d", verH, verL)
	return version, nil
}

func ClientgetValidDate() (validDate string, err error) {
	serial, err := ClientReadLicense()
	if err != nil {
		return serial, err
	}

	verDaySign, err := base64.StdEncoding.DecodeString(serial)
	if err != nil {
		return validDate, err
	}

	validDate = fmt.Sprintf("%s-%s-%s", string(verDaySign[4:8]), string(verDaySign[8:10]), string(verDaySign[10:12]))
	return validDate, nil
}
