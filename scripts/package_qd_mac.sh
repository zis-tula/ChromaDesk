#!/bin/bash

{
    cd "$(dirname "$0")"
    script_path=$(pwd)
    cd - > /dev/null
} &> /dev/null

old_cd=$(pwd)
cd "$(dirname "$0")"

build_mode=Release

while [ $# -gt 0 ]; do
    case "$(echo "$1" | tr '[:upper:]' '[:lower:]')" in
        debug)   build_mode=Debug ;;
        release) build_mode=Release ;;
    esac
    shift
done

publish_path="$script_path/../publish/$build_mode"
app_path="$publish_path/ChromaDesk.app"
dmg_path="$publish_path/ChromaDesk.dmg"

echo
echo
echo "---------------------------------------------------------------"
echo "create DMG"
echo "---------------------------------------------------------------"

if [ ! -d "$app_path" ]; then
    echo "[!] error: $app_path not found"
    echo "[!] please run publish_qd_mac.sh $build_mode first"
    cd "$old_cd"
    exit 1
fi

echo "[*] app path: $app_path"
echo "[*] dmg path: $dmg_path"

if [ -f "$dmg_path" ]; then
    rm -f "$dmg_path"
fi

tmp_dir=$(mktemp -d)
echo "[*] staging in: $tmp_dir"

cp -R "$app_path" "$tmp_dir/"
ln -s /Applications "$tmp_dir/Applications"

echo "[*] creating DMG..."
hdiutil create \
    -volname "ChromaDesk" \
    -srcfolder "$tmp_dir" \
    -ov \
    -format UDZO \
    "$dmg_path"

rm -rf "$tmp_dir"

if [ ! -f "$dmg_path" ]; then
    echo "[!] failed to create DMG"
    cd "$old_cd"
    exit 1
fi

echo
echo
echo "---------------------------------------------------------------"
echo "[*] DMG package finished!"
echo "---------------------------------------------------------------"
echo "[*] output: $dmg_path"
echo

cd "$old_cd"
exit 0
