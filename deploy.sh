#!/bin/bash

echo "Deploying..."
rsync -a --exclude='/.git' --exclude='.env' ../sesame ubuntu@niltouch.cn:~/go/src
echo "Done!"
