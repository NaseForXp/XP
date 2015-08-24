package ET99

import (
	"../serial"
	"errors"
	"fmt"
	"unsafe"
)

/*
#include <stdio.h>
#include <time.h>
#include "FT_ET99_API.h"

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
int ET_Et99_set_code(char *code){
	int count = 0;
	int err   = 0;
	ET_STATUS ret   = 0;
	unsigned char *pid      = NULL;
	unsigned char  pid1[10] = "FFFFFFFF";
	unsigned char  pid2[10] = "FFC5EB78";
	unsigned char  pin[20]  = "FFFFFFFFFFFFFFFF";
	char wbuf[60];

	ET_HANDLE hET99;

	// 打开动态库
	if ( ET_Dll_Open() != 0 ){
		return -1;
	}

	// 导出动态库
	if ( ET_Dll_Export() != 0 ){
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
			err = -2;
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
		err = -3;
		goto out2;
	}

	// 打开锁
	ret = et_OpenToken_Func(&hET99, pid, 1);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_OpenToken : %d\n", ret);
		err = -4;
		goto out2;
	}

	// 验证密码
	pin[16] = '\0';
	ret = et_Verify_Func(hET99, ET_VERIFY_SOPIN, pin);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Verify : %d\n\n", ret);
		err = -5;
		goto out;
	}

	// 设置PID
	if ( pid == pid1 ){
		ret = et_GenPID_Func(hET99, 6, "123456", pid2);
		if ( ret != ET_SUCCESS ){
			printf("Err : et_GenPID : %d\n\n", ret);
			err = -6;
			goto out;
		}
	}
	// 写数据
	strncpy(wbuf, code, (strlen(code) >= 50) ? 50 : strlen(code));
	ret = et_Write_Func(hET99, 0, 50, wbuf);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Write : %d\n", ret);
		err = -4;
		goto out;
	}

out:
	// 关闭锁
	ret = et_CloseToken_Func(hET99);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_CloseToken : %d\n\n", ret);
		err = -5;
	}

out2:
	// 释放动态库
	ET_Dll_Close();
	return err;
}

// 读取key内容
int ET_Et99_read_code(void *code, int codeLen){
	int count = 0;
	int err   = 0;
	ET_STATUS ret   = 0;
	unsigned char  pid[10] = "FFC5EB78";
	unsigned char  pin[20]  = "FFFFFFFFFFFFFFFF";
	char rbuf[60];

	ET_HANDLE hET99;

	// 打开动态库
	if ( ET_Dll_Open() != 0 ){
		return -1;
	}

	// 导出动态库
	if ( ET_Dll_Export() != 0 ){
		goto out;
	}

	// 检查key有没有插入
	pid[8] = '\0';
	ret = et_FindToken_Func(pid, &count);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_FindToken : %d\n", ret);
		err = -2;
		goto out2;
	}

	// 插入了多个key出错
	if ( count > 1 ){
		printf("Err : et_FindToken : 插入key太多 : %d\n", count);
		err = -3;
		goto out2;
	}

	// 打开锁
	ret = et_OpenToken_Func(&hET99, pid, 1);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_OpenToken : %d\n", ret);
		err = -4;
		goto out2;
	}

	// 验证密码
	pin[16] = '\0';
	ret = et_Verify_Func(hET99, ET_VERIFY_USERPIN, pin);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Verify : %d\n\n", ret);
		err = -5;
		goto out;
	}

	// 读数据
	ret = et_Read_Func(hET99, 0, 50, rbuf);
	if ( ret != ET_SUCCESS ){
		printf("Err : et_Read : %d\n", ret);
		err = -6;
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
		err = -5;
	}

out2:
	// 释放动态库
	ET_Dll_Close();
	return err;
}


*/
import "C"

// 验证硬件是否插入usbkey，以及key是否匹配
func Et99_check_login() (err error) {
	hcode, err := serial.ClientGetRegInfo()
	if err != nil {
		return err
	}

	var rbuf [50]C.char
	r := C.ET_Et99_read_code((unsafe.Pointer)(&rbuf), 50)
	if r != 0 {
		return errors.New(fmt.Sprintf("错误:ET_Et99_read_code() ret = %d", r))
	}

	rcode := C.GoString(&rbuf[0])
	if hcode == rcode {
		return nil
	}

	return errors.New("错误:请插入正确的USB_Key")
}

/*
func main() {
	code := "04291F39-C26B2600-C07A9F32"

	fmt.Println(code)

		// 一个key有写入次数上限
		r := C.ET_Et99_set_code(C.CString(code))
		if r != 0 {
			fmt.Println("写入失败")
		} else {
			fmt.Println("写入成功")
		}


	fmt.Println(Et99_check_login())
}
*/
