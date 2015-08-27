package et99

//package main

import (
	"errors"
	"fmt"
	"unsafe"

	"../serial"
)

/*
#include <stdio.h>
#include <time.h>
#include "FT_ET99_API.h"

typedef struct UserInfo_t{
	int   type;      // 0:客户端 1:管理中心
	char  info[20];  // 简单描述
}UserInfo;

#define KEY_TYPE_Center 1
#define KEY_TYPE_Client 2

// 导出函数列表
typedef ET_STATUS ET_API (*et_FindToken_t) (unsigned char* pid,int * count);
typedef ET_STATUS ET_API (*et_OpenToken_t) (ET_HANDLE* hHandle,unsigned char* pid,int index);
typedef ET_STATUS ET_API (*et_CloseToken_t)(ET_HANDLE hHandle);
typedef ET_STATUS ET_API (*et_Read_t)      (ET_HANDLE hHandle,WORD offset,int Len,unsigned char* pucReadBuf);
typedef ET_STATUS ET_API (*et_Write_t)     (ET_HANDLE hHandle,WORD offset,int Len,unsigned char* pucWriteBuf);
typedef ET_STATUS ET_API (*et_Verify_t)    (ET_HANDLE hHandle,int Flags,unsigned char* pucPIN);
typedef ET_STATUS ET_API (*et_GenPID_t)    (ET_HANDLE hHandle,int SeedLen,unsigned char* pucSeed,unsigned char* pid);

HINSTANCE gHdll;
et_FindToken_t       et_FindToken_Func;
et_OpenToken_t       et_OpenToken_Func;
et_CloseToken_t      et_CloseToken_Func;
et_Read_t            et_Read_Func;
et_Write_t           et_Write_Func;
et_Verify_t          et_Verify_Func;
et_GenPID_t          et_GenPID_Func;

int ET_Dll_Open(){
	gHdll = LoadLibrary("FT_ET99_API.dll");
	if( gHdll == NULL ){
		return -1;
	}
	return 0;
}

void ET_Dll_Close(){
	if ( gHdll != NULL ){
		FreeLibrary(gHdll);
	}
}

int ET_Dll_Export(){
	et_FindToken_Func = (et_FindToken_t)GetProcAddress(gHdll, "et_FindToken");
	if( et_FindToken_Func == NULL ){
		printf("Err : Export : et_FindToken\n");
		return -1;
	}

	et_OpenToken_Func = (et_OpenToken_t)GetProcAddress(gHdll, "et_OpenToken");
	if( et_OpenToken_Func == NULL ){
		printf("Err : Export : et_OpenToken\n");
		return -1;
	}

	et_CloseToken_Func = (et_CloseToken_t)GetProcAddress(gHdll, "et_CloseToken");
	if( et_CloseToken_Func == NULL ){
		printf("Err : Export : et_CloseToken\n");
		return -1;
	}

	et_Read_Func = (et_Read_t)GetProcAddress(gHdll, "et_Read");
	if( et_Read_Func == NULL ){
		printf("Err : Export : et_Read\n");
		return -1;
	}

	et_Write_Func = (et_Write_t)GetProcAddress(gHdll, "et_Write");
	if( et_Write_Func == NULL ){
		printf("Err : Export : et_Write\n");
		return -1;
	}

	et_Verify_Func = (et_Verify_t)GetProcAddress(gHdll, "et_Verify");
	if( et_Verify_Func == NULL ){
		printf("Err : Export : et_Verify\n");
		return -1;
	}

	et_GenPID_Func = (et_GenPID_t)GetProcAddress(gHdll, "et_GenPID");
	if( et_GenPID_Func == NULL ){
		printf("Err : Export : et_GenPID\n");
		return -1;
	}

	return 0;
}


// 初始化key
int ET_Et99_set_code(int type, char *info, char *code){
	int count = 0;
	int err   = 0;
	ET_STATUS ret   = 0;
	unsigned char *pid      = NULL;
	unsigned char  pid1[10] = "FFFFFFFF";
	unsigned char  pid2[10] = "FFC5EB78";
	unsigned char  pin[20]  = "FFFFFFFFFFFFFFFF";
	char wbuf[60];
	UserInfo uinfo;

	ET_HANDLE hET99;

	if ( type != KEY_TYPE_Center && type != KEY_TYPE_Client ){
		printf("KeyType Not Set.\n");
		return -1;
	}

	memset(&uinfo, 0x00, sizeof(UserInfo));
	uinfo.type = type;
	if ( strlen(info) >= 20 ){
		memcpy(uinfo.info, info, 20);
	}
	else{
		memcpy(uinfo.info, info, strlen(info));
	}


	// 打开动态库
	if ( ET_Dll_Open() != 0 ){
		printf("Can't Open FT_ET99_API.dll.\n");
		return -2;
	}

	// 导出动态库
	if ( ET_Dll_Export() != 0 ){
		err = -3;
		goto out;
	}

	// 检查key有没有插入
	pid1[8] = '\0';
	pid2[8] = '\0';
	ret = et_FindToken_Func(pid1, &count);
	if ( ret != ET_SUCCESS ){
		ret = et_FindToken_Func(pid2, &count);
		if ( ret != ET_SUCCESS ){
			printf("Err : et_FindToken : %d\n", ret);
			err = -4;
			goto out2;
		}
		else{
			pid = pid2;
		}
	}
	else{
		pid = pid1;
	}

	// 插入了多个key出错
	if ( count > 1 ){
		printf("Err : et_FindToken : 插入key太多 : %d\n", count);
		err = -5;
		goto out2;
	}

	// 打开锁
	ret = et_OpenToken_Func(&hET99, pid, 1);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_OpenToken : %d\n", ret);
		err = -6;
		goto out2;
	}

	// 验证密码
	pin[16] = '\0';
	ret = et_Verify_Func(hET99, ET_VERIFY_SOPIN, pin);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Verify : %d\n\n", ret);
		err = -7;
		goto out;
	}

	// 设置PID
	if ( pid == pid1 ){
		ret = et_GenPID_Func(hET99, 6, "123456", pid2);
		if ( ret != ET_SUCCESS ){
			printf("Err : et_GenPID : %d\n\n", ret);
			err = -8;
			goto out;
		}
	}

	// 写数据 前50字节 - 用户信息
	memcpy(wbuf, (unsigned char *)&uinfo, sizeof(UserInfo));
	ret = et_Write_Func(hET99, 0, 50, wbuf);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Write UserInfo: %d\n", ret);
		err = -9;
		goto out;
	}

	// 写数据 50~100字节 - 机器绑定信息
	memset(wbuf, 0x00, 60);
	strncpy(wbuf, code, (strlen(code) >= 50) ? 50 : strlen(code));
	ret = et_Write_Func(hET99, 50, 50, wbuf);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Write Machine Info: %d\n", ret);
		err = -10;
		goto out;
	}

out:
	// 关闭锁
	ret = et_CloseToken_Func(hET99);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_CloseToken : %d\n\n", ret);
		err = -11;
	}

out2:
	// 释放动态库
	ET_Dll_Close();
	return err;
}

// 读取key内容
int ET_Et99_read_code(void *code, int codeLen, int *type, void *info, int infoLen){
	int count = 0;
	int err   = 0;
	ET_STATUS ret   = 0;
	unsigned char  pid[10] = "FFC5EB78";
	unsigned char  pin[20]  = "FFFFFFFFFFFFFFFF";
	char rbuf[60];
	UserInfo *uinfo = (UserInfo *)rbuf;

	ET_HANDLE hET99;

	// 打开动态库
	if ( ET_Dll_Open() != 0 ){
		printf("Can't Open FT_ET99_API.dll.\n");
		return -1;
	}

	// 导出动态库
	if ( ET_Dll_Export() != 0 ){
		err = -2;
		goto out;
	}

	// 检查key有没有插入
	pid[8] = '\0';
	ret = et_FindToken_Func(pid, &count);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_FindToken : %d\n", ret);
		err = -3;
		goto out2;
	}

	// 插入了多个key出错
	if ( count > 1 ){
		printf("Err : et_FindToken : 插入key太多 : %d\n", count);
		err = -4;
		goto out2;
	}

	// 打开锁
	ret = et_OpenToken_Func(&hET99, pid, 1);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_OpenToken : %d\n", ret);
		err = -5;
		goto out2;
	}

	// 验证密码
	pin[16] = '\0';
	ret = et_Verify_Func(hET99, ET_VERIFY_USERPIN, pin);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Verify : %d\n\n", ret);
		err = -6;
		goto out;
	}

	// 读数据 - 前50字节用户信息
	memset(rbuf, 0x00, 60);
	ret = et_Read_Func(hET99, 0, 50, rbuf);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Read : %d\n", ret);
		err = -7;
		goto out;
	}
	(*type) = uinfo->type;
	if ( infoLen >= 20 ){
		memcpy(info, uinfo->info, 20);
	}
	else{
		memcpy(info, uinfo->info, infoLen);
	}

	// 读数据 - 前50~100字节硬件绑定信息
	ret = et_Read_Func(hET99, 50, 50, rbuf);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Read : %d\n", ret);
		err = -8;
		goto out;
	}

	if ( codeLen >= 50 ){
		memcpy(code, rbuf, 50);
	}
	else{
		memcpy(code, rbuf, codeLen);
	}

out:
	// 关闭锁
	ret = et_CloseToken_Func(hET99);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_CloseToken : %d\n\n", ret);
		err = -9;
	}

out2:
	// 释放动态库
	ET_Dll_Close();
	return err;
}


*/
import "C"

// 验证硬件是否插入usbkey，以及key是否匹配
// 客户端验证机器码 xxxxxxxx-xxxxxxxx-xxxxxxx
func Et99_check_client_login() (err error) {
	hcode, err := serial.ClientGetRegInfo()
	if err != nil {
		return err
	}

	var keytype C.int
	var info [50]C.char
	var mcode [50]C.char

	r := C.ET_Et99_read_code((unsafe.Pointer)(&mcode), 50, &keytype, (unsafe.Pointer)(&info), 50)
	switch r {
	case 0:
		break
	case -1:
		return errors.New("错误：没有找到FT_ET99_API.dll")
	case -2:
		return errors.New("错误：导出FT_ET99_API.dll失败")
	case -3:
		return errors.New("错误：查找USBKEY失败")
	case -4:
		return errors.New("错误：查找到多个USBKEY")
	case -5:
		return errors.New("错误：打开USBKEY失败")
	case -6:
		return errors.New("错误：验证USBKEY PIN密码失败")
	case -7:
		return errors.New("错误：读取USBKEY信息失败")
	case -8:
		return errors.New("错误：读取USBKEY绑定信息失败")
	case -9:
		return errors.New("错误：关闭USBKEY失败")
	default:
		return errors.New("错误：未知错误")
	}

	fmt.Println("Localcode: ", hcode)
	//fmt.Println("type: ", keytype, "info: ", C.GoString(&info[0]))
	fmt.Println("keycode: ", C.GoString(&mcode[0]))

	if keytype != C.KEY_TYPE_Client {
		return errors.New("错误:Key类型错误，请插入正确的USB_Key")
	}

	rcode := C.GoString(&mcode[0])
	if hcode == rcode {
		return nil
	}

	return errors.New("错误:Key绑定信息不匹配，请插入正确的USB_Key")
}

/*
func main() {
	code := "04291F39-C26B2600-17FA29EE"
	info := "用户信息"

	fmt.Println(code, info)

	// 一个key有写入次数上限
	//if C.ET_Et99_set_code(C.KEY_TYPE_Client, C.CString(info), C.CString(code)) != 0 {
	//	fmt.Println("写入失败")
	//} else {
	//	fmt.Println("写入成功")
	//}

	e := Et99_check_client_login()
	if e == nil {
		fmt.Println("USBKEY 验证成功")
	} else {
		fmt.Println(e.Error())
	}
}
*/
