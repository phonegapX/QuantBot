#!/bin/bash

mkdir ./cache
xgo --targets=windows/*,darwin/amd64,linux/amd64,linux/386,linux/arm --dest=cache ./

osarchs=(windows_amd64 windows_386 darwin_amd64 linux_amd64 linux_386 linux_arm)
files=(QuantBot-windows-4.0-amd64.exe QuantBot-windows-4.0-386.exe QuantBot-darwin-10.6-amd64 QuantBot-linux-amd64 QuantBot-linux-386 QuantBot-linux-arm-5)

for i in 0 1 2 3 4 5; do
  mkdir cache/QuantBot_${osarchs[${i}]}
  mkdir cache/QuantBot_${osarchs[${i}]}/web
  mkdir cache/QuantBot_${osarchs[${i}]}/custom
  cp LICENSE cache/QuantBot_${osarchs[${i}]}/LICENSE
  cp -r plugin cache/QuantBot_${osarchs[${i}]}/plugin
  cp README.md cache/QuantBot_${osarchs[${i}]}/README.md
  cp -r web/dist cache/QuantBot_${osarchs[${i}]}/web/dist
  cp custom/config.ini cache/QuantBot_${osarchs[${i}]}/custom/config.ini
  cp custom/config.ini cache/QuantBot_${osarchs[${i}]}/custom/config.default.ini
  cd cache
  if [ ${i} -lt 2 ]
  then
    mv ${files[${i}]} QuantBot_${osarchs[${i}]}/QuantBot.exe
    zip -r QuantBot_${osarchs[${i}]}.zip QuantBot_${osarchs[${i}]}
  else
    mv ${files[${i}]} QuantBot_${osarchs[${i}]}/QuantBot
    tar -zcvf QuantBot_${osarchs[${i}]}.tar.gz QuantBot_${osarchs[${i}]}
  fi
  rm -rf QuantBot_${osarchs[${i}]}
  cd ..
done
