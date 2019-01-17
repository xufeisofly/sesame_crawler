#!/bin/bash

echo "Deploy Start"
echo "Rsync files..."
rsync -a --exclude='.git' --exclude='.env' ../sesame ubuntu@niltouch.cn:~/go/src
echo "Done!"
echo "Running Local Script..."
ssh -t ubuntu@niltouch.cn 'cd go/src/sesame; bash -s' < local.sh
echo "Done!"
echo "Deploy Success"
