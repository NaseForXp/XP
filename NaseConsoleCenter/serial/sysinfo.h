#ifndef _SYSINFO_H_
#define _SYSINFO_H_

#include <Windows.h>
//#include <string>
//using namespace std;

#ifdef SYSINFO_EXPORTS
#define SYSINFO_API  __declspec(dllexport)
#else
#define SYSINFO_API  __declspec(dllimport)
#endif

#ifdef __cplusplus
extern "C" {
#endif

	SYSINFO_API	char* GetMotherboardInfo(char* buf, int len);

	SYSINFO_API char* GetPhyMemoryInfo(char* buf, int len);

	SYSINFO_API char* GetCpuInfo(char* buf, int len);

	SYSINFO_API char* GetBiosInfo(char* buf, int len);

	SYSINFO_API char* GetHardDiskInfo(char* buf, int len);
	SYSINFO_API int	  InitCom();

#ifdef __cplusplus
}
#endif

#endif  // _SEWINDOWS_H_
