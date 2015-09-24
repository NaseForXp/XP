@echo off

@"C:\\Program Files\\NaseForXP\\NaseClient\\instsrv.exe"  NaseForXPService "C:\\Program Files\\NaseForXP\\NaseClient\\srvany.exe"

@reg add HKLM\SYSTEM\CurrentControlSet\Services\NaseForXPService\Parameters /f

@reg add HKLM\SYSTEM\CurrentControlSet\Services\NaseForXPService\Parameters /v "Application"  /t REG_SZ /d "C:\\Program Files\\NaseForXP\\NaseClient\\NaseService.exe"

@reg add HKLM\SYSTEM\CurrentControlSet\Services\NaseForXPService\Parameters /v "AppDirectory"  /t REG_SZ /d "C:\\Program Files\\NaseForXP\\NaseClient"

@reg add HKLM\SYSTEM\CurrentControlSet\Services\NaseForXPService\Parameters /v "AppParameters"  /t REG_SZ /d ""

@sc config NaseForXPService start= auto

@net start NaseForXPService


