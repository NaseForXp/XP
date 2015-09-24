@echo off
@"C:\\Program Files\\NaseForXP\\NaseConsoleCenter\\instsrv.exe"  NaseForXPConsoleCenter "C:\\Program Files\\NaseForXP\\NaseConsoleCenter\\srvany.exe"

@reg add HKLM\SYSTEM\CurrentControlSet\Services\NaseForXPConsoleCenter\Parameters /f

@reg add HKLM\SYSTEM\CurrentControlSet\Services\NaseForXPConsoleCenter\Parameters /v "Application"  /t REG_SZ /d "C:\\Program Files\\NaseForXP\\NaseConsoleCenter\\NaseConsoleCenter.exe"

@reg add HKLM\SYSTEM\CurrentControlSet\Services\NaseForXPConsoleCenter\Parameters /v "AppDirectory"  /t REG_SZ /d "C:\\Program Files\\NaseForXP\\NaseConsoleCenter"

@reg add HKLM\SYSTEM\CurrentControlSet\Services\NaseForXPConsoleCenter\Parameters /v "AppParameters"  /t REG_SZ /d ""

@sc config NaseForXPConsoleCenter start= auto

@net start NaseForXPConsoleCenter


