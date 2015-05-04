#ifndef _GO_SEWINDOWS_H_
#define _GO_SEWINDOWS_H_

int C_SewinInit();
BOOLEAN C_SewinSetModeNotify();
BOOLEAN C_SewinSetModeIntercept();
BOOLEAN C_SewinRegOps();

BOOLEAN C_file_create(WCHAR *user_name, WCHAR *process, WCHAR *file_path);

#endif
