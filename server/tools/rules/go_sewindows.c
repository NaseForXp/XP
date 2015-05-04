#include <stdio.h>
#include "sewindows.h"

#pragma comment(lib,"Advapi32.lib") 
#pragma comment(lib,"User32.lib") 


typedef BOOLEAN(*fsewin_init)();
typedef BOOLEAN(*fsewin_setoption)(int mode, int type);
typedef BOOLEAN(*fsewin_register_opt)(struct sewin_operations *ops);

fsewin_init         monitor_sewin_init;
fsewin_setoption    monitor_sewin_setoption;
fsewin_register_opt monitor_sewin_register_opt;

char *WcharToChar(WCHAR * WStr){
	int   len = WideCharToMultiByte(CP_ACP, 0, WStr, wcslen(WStr), NULL, 0, NULL, NULL);   
    char *m_char = (char *)malloc(len + 1);
	
	if (m_char == NULL){
		return NULL;
	}
	
    WideCharToMultiByte(CP_ACP, 0, WStr, wcslen(WStr), m_char, len, NULL, NULL);   
    m_char[len] = '\0';   
	
    return m_char;
}

BOOLEAN C_file_create(WCHAR *user_name, WCHAR *process, WCHAR *file_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(file_path);
	
	ret = Go_file_create(uname, proc, fpath);
	if(uname != NULL){
		free(uname);	
	}
	if(proc != NULL){
		free(proc);	
	}
	if(fpath != NULL){
		free(fpath);	
	}
	return ret;
}

int C_SewinInit(){
	int   ret   = 0;
    HMODULE handle;

    // step1. loadLibrary sewindows.dll
    handle                     = LoadLibrary("sewindows.dll");
    monitor_sewin_init         = (fsewin_init)GetProcAddress(handle, "sewin_init");
    monitor_sewin_setoption    = (fsewin_setoption)GetProcAddress(handle, "sewin_setoption");
    monitor_sewin_register_opt = (fsewin_register_opt)GetProcAddress(handle, "sewin_register_opt");

    if (monitor_sewin_init == NULL || monitor_sewin_setoption == NULL || monitor_sewin_register_opt == NULL)
    {
        return -1;
    }

    // step2. init sewindows
    BOOLEAN bret = monitor_sewin_init();
    if ( !bret )
    {
        return -2;
    }
	
	return 0;
}

BOOLEAN C_SewinSetModeNotify(){	
	return monitor_sewin_setoption(SEWIN_MODE_NOTIFY, SEWIN_TYPE_SCVDRV | SEWIN_TYPE_FILE | SEWIN_TYPE_PROC | SEWIN_TYPE_REG);
}

BOOLEAN C_SewinSetModeIntercept(){	
	return monitor_sewin_setoption(SEWIN_MODE_INTERCEPT, SEWIN_TYPE_SCVDRV | SEWIN_TYPE_FILE | SEWIN_TYPE_PROC | SEWIN_TYPE_REG);
}

BOOLEAN C_SewinRegOps(){
	struct sewin_operations ops;
	memset(&ops, 0x00, sizeof(struct sewin_operations));
	ops.file_create = C_file_create;
    return monitor_sewin_register_opt(&ops);
}
