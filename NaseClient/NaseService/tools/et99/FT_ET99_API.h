/*
[]=========================================================================[]

Copyright(C) 2006, Feitian Technologies Co., Ltd.
All rights reserved.
 
FILE:
	FT_ET99_API.h
	
DESC:
	et99 interface and defines 
		
[]=========================================================================[]
*/

#ifndef  __FT_ET99_HEADER_H
#define  __FT_ET99_HEADER_H

#include <windows.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef ET_STATUS
#define ET_STATUS			unsigned long
#define ET_API				__stdcall
typedef void*				ET_HANDLE;
#endif

//Return Value Status Definition
#define	  ET_SUCCESS				0x00			//函数执行成功
#define   ET_ACCESS_DENY            0x01			//访问被拒绝，权限不够
#define   ET_COMMUNICATIONS_ERROR   0x02			//通讯错误，没有打开设备
#define   ET_INVALID_PARAMETER      0x03			//无效的参数，参数出错
#define   ET_NOT_SET_PID            0x04			//没有设置PID
#define   ET_UNIT_NOT_FOUND 	    0x05			//打开指定的设备失败
#define   ET_HARD_ERROR             0x06			//硬件错误
#define   ET_UNKNOWN                0x07			//一般性错误
#define	  ET_PIN_ERR_MASK			0x0F			//验证PIN码掩码
#define	  ET_PIN_ERR_MAX		    0xFF            //验证PIN码最大剩余次数，如果验证错误还还回ET_PIN_ERR_MAX则表示PIN码永远不被锁死。

//define pin Flags
#define ET_VERIFY_USERPIN			0
#define ET_VERIFY_SOPIN				1

//define read/write flag
#define ET_USER_WRITE_READ				0
#define ET_USER_READ_ONLY				1

//出厂设置的PID
#define CONST_PID                   "FFFFFFFF"

ET_STATUS ET_API	et_FindToken(unsigned char* pid/*[in]*/,int * count);

ET_STATUS ET_API	et_OpenToken(ET_HANDLE* hHandle,unsigned char* pid,int index);

ET_STATUS ET_API	et_CloseToken(ET_HANDLE hHandle);

ET_STATUS ET_API	et_Read(ET_HANDLE hHandle,WORD offset,int Len,unsigned char* pucReadBuf);

ET_STATUS ET_API	et_Write(ET_HANDLE hHandle,WORD offset,int Len,unsigned char* pucWriteBuf);

ET_STATUS ET_API	et_GenPID(ET_HANDLE hHandle,int SeedLen,unsigned char* pucSeed,unsigned char* pid);

ET_STATUS ET_API	et_GenRandom(ET_HANDLE hHandle,unsigned char* pucRandBuf);

ET_STATUS ET_API	et_GenSOPIN(ET_HANDLE hHandle,int SeedLen,unsigned char* pucSeed, unsigned char* pucNewSoPIN);

ET_STATUS ET_API	et_ResetPIN(ET_HANDLE hHandle,unsigned char* pucSoPIN);

ET_STATUS ET_API	et_SetKey(ET_HANDLE hHandle,int Keyid, unsigned char* pucKeyBuf);

ET_STATUS ET_API	et_HMAC_MD5(ET_HANDLE hHandle,int keyID,int textLen,unsigned char* pucText,unsigned char *digest);

ET_STATUS ET_API	et_Verify(ET_HANDLE hHandle,int Flags,unsigned char* pucPIN);

ET_STATUS ET_API	et_ChangeUserPIN(ET_HANDLE hHandle,unsigned char* pucOldPIN,unsigned char* pucNewPIN);

ET_STATUS ET_API	et_ResetSecurityState(ET_HANDLE hHandle);

ET_STATUS ET_API	et_GetSN(ET_HANDLE hHandle,unsigned char* pucSN);

ET_STATUS ET_API	et_SetupToken(ET_HANDLE hHandle,BYTE bSoPINRetries,BYTE bUserPINRetries,BYTE bUserReadOnly,BYTE bBack);

ET_STATUS ET_API	et_TurnOnLED(ET_HANDLE hHandle);

ET_STATUS ET_API	et_TurnOffLED(ET_HANDLE hHandle);





/*{{{ MD5 HMAC Wrapper functions.*/

ET_STATUS ET_API
MD5_HMAC(unsigned char * pucText,        /* pointer to data stream        */
		 unsigned long   ulText_Len,     /* length of data stream         */
		 unsigned char * pucKey,         /* pointer to authentication key */
		 unsigned long   ulKey_Len,      /* length of authentication key  */
		 unsigned char * pucToenKey,     /* inner hashed key and  outer hashed key */
		 unsigned char * pucDigest );    /* caller digest to be filled in */
/*}}}*/

#ifdef __cplusplus 
}
#endif

#endif	//__FT_ET99_HEADER_H

