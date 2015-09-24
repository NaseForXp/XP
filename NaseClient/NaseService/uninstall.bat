@echo off

cd %ProgramFiles%\\NaseForXP\\NaseClient

taskkill /F /im NaseService.exe | findstr "成功" && (
	taskkill /F /im NASE客户端.exe

	@net stop NaseForXPService
	instsrv.exe NaseForXPService remove

	rd /S /Q .
	del %0
	echo "卸载成功"
) || (
	echo "卸载失败"
)
