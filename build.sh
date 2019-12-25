#!/usr/bin/env bash
#  sh /Users/kdaxrobot/Library/Preferences/GoLand2019.1/scratches/build.sh okex-future

cd $(pwd)

appname=$1
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" ./App.go
rm $appname -f
mv App $appname

upx $appname

workdir=$(cd $(dirname $0); pwd)

sh /Users/kdaxrobot/Library/Preferences/GoLand2019.1/scratches/send.sh $appname $workdir