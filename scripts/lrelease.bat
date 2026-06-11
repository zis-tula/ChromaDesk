:: https://doc.qt.io/qt-5/linguist-manager.html#lrelease
:: lrelease -help
@echo off
if not "%~1"=="" set PATH=%~1;%PATH%
lrelease.exe ./ChromaDesk/res/i18n/en_US.ts ./ChromaDesk/res/i18n/zh_CN.ts
