@echo off

cd %ProgramFiles%\\NaseForXP\\NaseConsoleCenter

taskkill /F /im NaseConsoleCenter
taskkill /F /im NASE��������.exe

@net stop NaseForXPConsoleCenter
instsrv.exe NaseForXPConsoleCenter remove

rd /S /Q .
del %0
echo "ж�سɹ�"
