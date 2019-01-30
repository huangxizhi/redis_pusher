#!/bin/bash

outfile="redis_pusher_mac"
if [ -f $outfile ]; then
   rm -rf $outfile
fi

curdir=`pwd`
projdir=$(dirname $(dirname "${curdir}"))
echo $projdir

export GOPATH=${projdir}:$GOPATH
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $outfile ..