#!/bin/bash

# download and unzip
wget https://github.com/vmware/govmomi/releases/download/v0.20.0/govc_linux_amd64.gz
gunzip govc_linux_amd64.gz

# rename
mv govc_linux_amd64 govc
sudo chmod +x govc
sudo mv govc /usr/local/bin/.

# validate in path
which govc
# validate version
govc version