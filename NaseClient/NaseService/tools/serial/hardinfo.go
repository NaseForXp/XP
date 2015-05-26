package serial

/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <windows.h>
#include <tchar.h>

static DWORD g_eax;   // 存储返回的eax
static DWORD g_ebx;   // 存储返回的ebx
static DWORD g_ecx;   // 存储返回的ecx
static DWORD g_edx;   // 存储返回的edx

void Executecpuid(DWORD veax)
{
       asm("cpuid"
           :"=a"(g_eax),
           "=b"(g_ebx),
           "=c"(g_ecx),
           "=d"(g_edx)
           :"a"(g_eax));
}
int isSupport;
int GetSerialNumber(WORD nibble[6])
{
        Executecpuid(1);                  // 执行cpuid，参数为 eax = 1
        isSupport = g_edx & (1<<18); // edx是否为1代表CPU是否存在序列号
        if (FALSE == isSupport)           // 不支持，返回false
        {
                return -1;
        }
        Executecpuid(3);              // 执行cpuid，参数为 eax = 3
        memcpy(&nibble[4], &g_eax, 4); // eax为最高位的两个WORD
        memcpy(&nibble[0], &g_ecx, 8); // ecx 和 edx为低位的4个WORD
 		return 0;
}

char * GetSn(char *buf)
{
	WORD nibble[8];
	int i = 0;

	memset(nibble, 0x00, sizeof(WORD) * 8);
	i = GetSerialNumber(nibble);

	memset(buf, 0x00, 64);
	sprintf(buf, "%04X%04X%04X%04X%04X%04X", nibble[0], nibble[1], nibble[2], nibble[3], nibble[4], nibble[5]);

	return buf;
}

*/
import "C"

type HardWareInfo struct {
	SerialNumber string
	Uuid         string
	CpuId        string
}

func GetSysInfo() (info HardWareInfo, err error) {
	var buf [64]C.char

	C.GetSn((*C.char)(&buf[0]))

	info.SerialNumber = ""
	info.Uuid = "lzf:24324962@qq.com"
	info.CpuId = C.GoString((*C.char)(&buf[0]))

	// 测试
	info.CpuId = "This is a test"
	return info, nil
}
