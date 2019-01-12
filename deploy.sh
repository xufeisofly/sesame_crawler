#!/bin/bash

echo "Deploying start"
echo "Rsync files..."
rsync -a --exclude='/.git' --exclude='.env' --exclude='/vendor' ../sesame ubuntu@niltouch.cn:~/go/src
echo "Done!"
echo "Running Local Script..."
ssh -t ubuntu@niltouch.cn 'cd go/src; bash -s' < local.sh
echo "Done!"
echo "Deploy Success!"
