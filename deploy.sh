#!/bin/bash

echo "Deploying..."
rsync -a --exclude='/.git' ../sesame ubuntu@118.89.236.14:~/go/src
echo "Done!"
