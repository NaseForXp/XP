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

BOOLEAN C_file_unlink(WCHAR *user_name, WCHAR *process, WCHAR *file_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(file_path);
	
	ret = Go_file_unlink(uname, proc, fpath);
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

BOOLEAN C_file_read(WCHAR *user_name, WCHAR *process, WCHAR *file_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(file_path);
	
	ret = Go_file_read(uname, proc, fpath);
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

BOOLEAN C_file_write(WCHAR *user_name, WCHAR *process, WCHAR *file_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(file_path);
	
	ret = Go_file_write(uname, proc, fpath);
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

BOOLEAN C_file_rename(WCHAR *user_name, WCHAR *process, WCHAR *src_file, WCHAR *new_name){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	char *newpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(src_file);
	newpath = WcharToChar(new_name);
	
	ret = Go_file_rename(uname, proc, fpath, newpath);
	if(uname != NULL){
		free(uname);	
	}
	if(proc != NULL){
		free(proc);	
	}
	if(fpath != NULL){
		free(fpath);	
	}
	if(newpath != NULL){
		free(newpath);	
	}
	return ret;
}

BOOLEAN C_dir_create(WCHAR *user_name, WCHAR *process, WCHAR *file_path){
	return C_file_create(user_name, process, file_path);
}

BOOLEAN C_dir_unlink(WCHAR *user_name, WCHAR *process, WCHAR *file_path){
	return C_file_unlink(user_name, process, file_path);
}

BOOLEAN C_dir_rename(WCHAR *user_name, WCHAR *process, WCHAR *src_file, WCHAR *new_name){
	return C_file_rename(user_name, process, src_file, new_name);
}

BOOLEAN C_disk_read(WCHAR *user_name, WCHAR *process, WCHAR *dir_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(dir_path);
	
	ret = Go_disk_read(uname, proc, fpath);
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

BOOLEAN C_disk_write(WCHAR *user_name, WCHAR *process, WCHAR *dir_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(dir_path);
	
	ret = Go_disk_write(uname, proc, fpath);
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

BOOLEAN C_disk_format(WCHAR *user_name, WCHAR *process, WCHAR *dir_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(dir_path);
	
	ret = Go_disk_formate(uname, proc, fpath);
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

BOOLEAN C_process_kill(WCHAR *user_name, WCHAR *process, WCHAR *dst_proc){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(dst_proc);
	
	ret = Go_process_kill(uname, proc, fpath);
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


BOOLEAN C_process_create_thread(WCHAR *user_name, WCHAR *process, WCHAR *dst_proc){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	fpath = WcharToChar(dst_proc);
	
	ret = Go_process_create_thread(uname, proc, fpath);
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

BOOLEAN C_service_create(WCHAR *user_name, WCHAR *process, WCHAR *service_name, WCHAR *bin_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *name  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	name = WcharToChar(service_name);
	fpath = WcharToChar(bin_path);
	
	ret = Go_service_create(uname, proc, name, fpath);
	if(uname != NULL){
		free(uname);	
	}
	if(proc != NULL){
		free(proc);	
	}
	if(name != NULL){
		free(name);	
	}
	if(fpath != NULL){
		free(fpath);	
	}
	return ret;
}

BOOLEAN C_service_delete(WCHAR *user_name, WCHAR *process, WCHAR *service_name){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *name  = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	name = WcharToChar(service_name);
	
	ret = Go_service_delete(uname, proc, name);
	if(uname != NULL){
		free(uname);	
	}
	if(proc != NULL){
		free(proc);	
	}
	if(name != NULL){
		free(name);	
	}
	return ret;
}

BOOLEAN C_service_change(WCHAR *user_name, WCHAR *process, WCHAR *service_name){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *name  = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	name = WcharToChar(service_name);
	
	ret = Go_service_change(uname, proc, name);
	if(uname != NULL){
		free(uname);	
	}
	if(proc != NULL){
		free(proc);	
	}
	if(name != NULL){
		free(name);	
	}
	return ret;
}

BOOLEAN C_driver_load(WCHAR *user_name, WCHAR *process, WCHAR *service_name, WCHAR *bin_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *name  = NULL;
	char *fpath = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	name = WcharToChar(service_name);
	fpath = WcharToChar(bin_path);
	
	ret = Go_driver_load(uname, proc, name, fpath);
	if(uname != NULL){
		free(uname);	
	}
	if(proc != NULL){
		free(proc);	
	}
	if(name != NULL){
		free(name);	
	}
	if(fpath != NULL){
		free(fpath);	
	}
	return ret;
}

BOOLEAN C_reg_set_value(WCHAR *user_name, WCHAR *process, WCHAR *reg_path, WCHAR *reg_value){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *rpath  = NULL;
	char *rvalue = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	rpath = WcharToChar(reg_path);
	rvalue = WcharToChar(reg_value);
	
	ret = Go_reg_set_value(uname, proc, rpath, rvalue);
	if(uname != NULL){
		free(uname);	
	}
	if(proc != NULL){
		free(proc);	
	}
	if(rpath != NULL){
		free(rpath);	
	}
	if(rvalue != NULL){
		free(rvalue);	
	}
	return ret;
}
/*
BOOLEAN C_reg_create_key(WCHAR *user_name, WCHAR *process, WCHAR *reg_path){
	BOOLEAN ret = FALSE;
	char *uname = NULL;
	char *proc  = NULL;
	char *rpath  = NULL;
	
	uname = WcharToChar(user_name);
	proc = WcharToChar(process);
	rpath = WcharToChar(reg_path);
	
	ret = Go_reg_create_key(uname, proc, rpath);
	if(uname != NULL){
		free(uname);	
	}
	if(proc != NULL){
		free(proc);	
	}
	if(rpath != NULL){
		free(rpath);	
	}
	return ret;
}
*/
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
	ops.file_unlink = C_file_unlink;
	ops.file_read   = C_file_read;
	ops.file_write  = C_file_write;
	ops.file_rename = C_file_rename;
	
	ops.dir_create  = C_dir_create;
	ops.dir_unlink  = C_dir_unlink;
	ops.dir_rename  = C_dir_rename;	
	
	ops.process_create_thread  = C_process_create_thread;
	ops.process_kill  = C_process_kill;
	
	ops.disk_read   = C_disk_read;
	ops.disk_write  = C_disk_write;
	ops.disk_format = C_disk_format;
	
	ops.service_create = C_service_create;
	ops.service_delete = C_service_delete;
	ops.service_change = C_service_change;
	ops.driver_load  = C_driver_load;
	
	ops.reg_set_value  = C_reg_set_value;
	//ops.reg_create_key = C_reg_create_key;
	
    return monitor_sewin_register_opt(&ops);
}


