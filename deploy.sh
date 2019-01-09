#!/bin/bash

echo "Deploying..."
rsync -a --exclude='/.git' ../sesame ubuntu@niltouch.cn:~/go/src
echo "Done!"
