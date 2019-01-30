#!/bin/bash

outfile="redis_pusher.exe"
if [ -f $outfile ]; then
   rm -rf $outfile
fi

curdir=`pwd`
projdir=$(dirname $(dirname "${curdir}"))
echo $projdir

export GOPATH=${projdir}:$GOPATH
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $outfile ..
