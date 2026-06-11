# https://doc.qt.io/qt-5/linguist-manager.html#lrelease
# lrelease -help
if [ -n "$1" ]; then
    export PATH="$1:$PATH"
fi
lrelease ./ChromaDesk/res/i18n/en_US.ts ./ChromaDesk/res/i18n/zh_CN.ts
