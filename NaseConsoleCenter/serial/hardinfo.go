package serial

/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <windows.h>
#include <tchar.h>
#include "sysinfo.h"

// 导出函数列表
typedef int	   (*InitCom_t)();
typedef char * (*GetMotherboardInfo_t) (char* buf,int len);
typedef char * (*GetPhyMemoryInfo_t)   (char* buf,int len);
typedef char * (*GetCpuInfo_t)         (char* buf,int len);
typedef char * (*GetBiosInfo_t)        (char* buf,int len);
typedef char * (*GetHardDiskInfo_t)    (char* buf,int len);

HINSTANCE gHdll;
InitCom_t                  InitCom_Func;
GetMotherboardInfo_t       GetMotherboardInfo_Func;
GetPhyMemoryInfo_t         GetPhyMemoryInfo_Func;
GetCpuInfo_t               GetCpuInfo_Func;
GetBiosInfo_t              GetBiosInfo_Func;
GetHardDiskInfo_t          GetHardDiskInfo_Func;

int SYS_Dll_Open(){
	gHdll = LoadLibrary("SysInfo.dll");
	if( gHdll == NULL ){
		return -1;
	}
	return 0;
}

void SYS_Dll_Close(){
	if ( gHdll != NULL ){
		FreeLibrary(gHdll);
	}
}

int SYS_Dll_Export(){
	InitCom_Func = (InitCom_t)GetProcAddress(gHdll, "InitCom");
	if( InitCom_Func == NULL ){
		printf("Err : Export : InitCom\n");
		return -1;
	}

	GetMotherboardInfo_Func = (GetMotherboardInfo_t)GetProcAddress(gHdll, "GetMotherboardInfo");
	if( GetMotherboardInfo_Func == NULL ){
		printf("Err : Export : GetMotherboardInfo\n");
		return -1;
	}

	GetPhyMemoryInfo_Func = (GetPhyMemoryInfo_t)GetProcAddress(gHdll, "GetPhyMemoryInfo");
	if( GetPhyMemoryInfo_Func == NULL ){
		printf("Err : Export : GetPhyMemoryInfo\n");
		return -1;
	}

	GetCpuInfo_Func = (GetCpuInfo_t)GetProcAddress(gHdll, "GetCpuInfo");
	if( GetCpuInfo_Func == NULL ){
		printf("Err : Export : GetCpuInfo\n");
		return -1;
	}

	GetBiosInfo_Func = (GetBiosInfo_t)GetProcAddress(gHdll, "GetBiosInfo");
	if( GetBiosInfo_Func == NULL ){
		printf("Err : Export : GetBiosInfo\n");
		return -1;
	}

	GetHardDiskInfo_Func = (GetHardDiskInfo_t)GetProcAddress(gHdll, "GetHardDiskInfo");
	if( GetHardDiskInfo_Func == NULL ){
		printf("Err : Export : GetMotherboardInfo\n");
		return -1;
	}

	return 0;
}


int GetCpuAndDiskInfo(char *cpuinfo, int cpuinfosize, char *diskinfo, int diskinfosize){
	int ret = 0;
	char buf[1024];
	ret = SYS_Dll_Open();
	if ( ret != 0 ){
		return -1;
	}

	ret = SYS_Dll_Export();
	if ( ret != 0 ){
		SYS_Dll_Close();
		return -2;
	}

	ret = InitCom_Func();
	if ( ret != 1 ){
		printf("Err:InitCom failed.\n");
		SYS_Dll_Close();
		return -3;
	}


	memset(buf, 0x00, 1024);
	memset(cpuinfo, 0x00, cpuinfosize);
	if( GetCpuInfo_Func(buf, 1024) == NULL ){
		printf("Err:GetCpuInfo \n");
	}
	strncpy(cpuinfo, buf, cpuinfosize -1);

	memset(buf, 0x00, 1024);
	memset(diskinfo, 0x00, diskinfosize);
	if( GetHardDiskInfo_Func(buf, 1024) == NULL ){
		printf("Err:GetHardDiskInfo \n");
	}
	strncpy(diskinfo, buf, cpuinfosize -1);

	SYS_Dll_Close();
	return 0;
}
*/
import "C"
import "hash/crc32"

type HardWareInfo struct {
	StaticInfo string // 静态信息
	CpuInfo    string // CPU信息
	DiskInfo   string // 硬盘信息
}

var GIsInit int = 0
var GlobalHardWareInfo HardWareInfo

// 获取CRC32
func GetCrc32(data []byte) (crc uint32) {
	h := crc32.NewIEEE()
	h.Write(data)
	return h.Sum32()
}

func GetSysInfo() (info HardWareInfo, err error) {
	if GIsInit == 1 {
		info = GlobalHardWareInfo
		return info, nil
	}

	var cpuinfo [1024]C.char
	var diskinfo [1024]C.char

	C.GetCpuAndDiskInfo((*C.char)(&cpuinfo[0]), 1024, (*C.char)(&diskinfo[0]), 1024)

	info.StaticInfo = "lzf:24324962@qq.com"
	info.CpuInfo = C.GoString((*C.char)(&cpuinfo[0]))
	info.DiskInfo = C.GoString((*C.char)(&diskinfo[0]))

	// 测试 - 正式版注释掉下一行
	//info.CpuId = "This is a test"
	GlobalHardWareInfo = info
	GIsInit = 1
	return info, nil
}
