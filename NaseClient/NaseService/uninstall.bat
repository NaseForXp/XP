@echo off

cd %ProgramFiles%\\NaseForXP\\NaseClient

taskkill /F /im NaseService.exe | findstr "�ɹ�" && (
	taskkill /F /im NASE�ͻ���.exe

	@net stop NaseForXPService
	instsrv.exe NaseForXPService remove

	rd /S /Q .
	del %0
	echo "ж�سɹ�"
) || (
	echo "ж��ʧ��"
)
