#ifndef _GO_ACCOUNT_H__
#define _GO_ACCOUNT_H__

typedef struct _security_set
{
	int  MinimumPasswordAge;
	int  MaximumPasswordAge;
	int  MinimumPasswordLength;
	int  PasswordComplexity;
	int  PasswordHistorySize;
	int  LockoutBadCount;
	int  LockoutDuration;
}security_set, *psecurity_set;

int get_account_security_set(psecurity_set pSet);
int set_account_security_set(psecurity_set pSet);

#endif
