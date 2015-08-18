#include <stdio.h>
#include <tchar.h>
#include <Windows.h>

#include "go_account.h"

int get_account_security_set(psecurity_set pSet)
{
	WIN32_FIND_DATAA	FindFileData;
	int					TryCount = 0;
	char				dir[MAX_PATH + 1] = { 0 };
	char				path[MAX_PATH + 1] = { 0 };
	char				cmd[MAX_PATH * 2 + 1] = { 0 };

	if (pSet == NULL)
	{
		return 0;
	}

	//GetModuleFileNameA(NULL, dir, MAX_PATH);
	//(strrchr(dir, '\\') + 1)[0] = 0; 	
	GetEnvironmentVariableA("TEMP", dir, MAX_PATH);
	
	sprintf(path, "%s%s", dir, "sec.ini");
	sprintf(cmd, "secedit /export /CFG \"%s\" /quiet", path);
	//system("secedit /export /CFG sec.ini /quiet");
	system(cmd);

	pSet->MaximumPasswordAge = GetPrivateProfileIntA("System Access", "MaximumPasswordAge", 0, path);
	pSet->MinimumPasswordAge = GetPrivateProfileIntA("System Access", "MinimumPasswordAge", 42, path);
	pSet->MinimumPasswordLength = GetPrivateProfileIntA("System Access", "MinimumPasswordLength", 0, path);
	pSet->PasswordComplexity = GetPrivateProfileIntA("System Access", "PasswordComplexity", 0, path);
	pSet->PasswordHistorySize = GetPrivateProfileIntA("System Access", "PasswordHistorySize", 0, path);
	pSet->LockoutBadCount = GetPrivateProfileIntA("System Access", "LockoutBadCount", 0, path);
	pSet->LockoutDuration = GetPrivateProfileIntA("System Access", "LockoutDuration", 30, path);
	DeleteFileA(path);
	return 1;
}

int set_account_security_set(psecurity_set pSet)
{
	WIN32_FIND_DATAA	FindFileData;
	int					TryCount = 0;
	char				dir[MAX_PATH + 1] = { 0 };
	char				path[MAX_PATH + 1] = { 0 };
	char				pathsdb[MAX_PATH + 1] = { 0 };
	char				cmd[MAX_PATH * 3 + 1] = { 0 };
	char buffer[65];
	if (pSet == NULL)
	{
		return 0;
	}
	
	//GetModuleFileNameA(NULL, dir, MAX_PATH);
	//(strrchr(dir, '\\') + 1)[0] = 0; 
	GetEnvironmentVariableA("TEMP", dir, MAX_PATH);
	
	sprintf(pathsdb, "%s%s", dir, "sec.sdb");
	sprintf(path, "%s%s", dir, "sec.ini");
		
	//system("secedit /export /CFG sec.ini /quiet");
	sprintf(cmd, "secedit /export /CFG \"%s\" /quiet", path);
	system(cmd);
	
	_itoa(pSet->MaximumPasswordAge, buffer, 10);
	WritePrivateProfileStringA("System Access", "MaximumPasswordAge", buffer, path);

	_itoa(pSet->MinimumPasswordAge, buffer, 10);
	WritePrivateProfileStringA("System Access", "MinimumPasswordAge", buffer, path);

	_itoa(pSet->MinimumPasswordLength, buffer, 10);
	WritePrivateProfileStringA("System Access", "MinimumPasswordLength", buffer, path);

	_itoa(pSet->PasswordComplexity, buffer, 10);
	WritePrivateProfileStringA("System Access", "PasswordComplexity", buffer, path);

	_itoa(pSet->PasswordHistorySize, buffer, 10);
	WritePrivateProfileStringA("System Access", "PasswordHistorySize", buffer, path);
	
	_itoa(pSet->LockoutBadCount, buffer, 10);
	WritePrivateProfileStringA("System Access", "LockoutBadCount", buffer, path);
	
	_itoa(pSet->LockoutDuration, buffer, 10);
	WritePrivateProfileStringA("System Access", "LockoutDuration", buffer, path);
	WritePrivateProfileStringA("System Access", "ResetLockoutCount", buffer, path);
	

	memset(cmd, 0x00, 3 * MAX_PATH)	;
	sprintf(cmd, "secedit /analyze /db \"%s\" /cfg \"%s\" /quiet", pathsdb, path);
	system(cmd);
	
	memset(cmd, 0x00, 3 * MAX_PATH)	;
	sprintf(cmd, "secedit /configure /db \"%s\"  /quiet", pathsdb);
	system(cmd);
	
	//system("secedit /analyze /db  sec.sdb  /cfg sec.ini /quiet");
	//system("secedit /configure /db sec.sdb  /quiet");

	DeleteFileA(path);
	DeleteFileA(pathsdb);
	return 1;
}

/*
int _tmain(int argc, TCHAR * argv[])
{
	security_set pSet;
	RtlZeroMemory(&pSet, sizeof(pSet));
	get_account_security_set(&pSet);
	pSet.PasswordComplexity  = 0;
	pSet.MaximumPasswordAge +=1;
	pSet.MinimumPasswordLength +=1;
	pSet.PasswordHistorySize +=1;
	set_account_security_set(&pSet);
	return 0;
}*/