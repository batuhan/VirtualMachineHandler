#!/bin/bash

# Ubuntu/Debian repositories have an older version that's not compatible 
# with our current configurations. Make sure to have version > 1.5.0

# download and unzip
wget https://releases.hashicorp.com/packer/1.5.1/packer_1.5.1_linux_amd64.zip
unzip packer_1.5.1_linux_amd64.zip -d packer
rm -rf packer_1.5.1_linux_amd64.zip


# move

DIR=$HOME/.packer.d
[ -d $DIR ] || mkdir -p $DIR
mv packer $DIR/

# make sure it's in the path
export PATH="$PATH:$HOME/.packer.d/packer"

# install vsphere plugins
wget https://github.com/jetbrains-infra/packer-builder-vsphere/releases/download/v2.3/packer-builder-vsphere-iso.linux
wget https://github.com/jetbrains-infra/packer-builder-vsphere/releases/download/v2.3/packer-builder-vsphere-clone.linux

# set executable
chmod +x packer-builder-vsphere-iso.linux
chmod +x packer-builder-vsphere-clone.linux

DIRP=$HOME/.packer.d/plugins
[ -d $DIRP ] || mkdir -p $DIRP
mv packer-builder-vsphere-*.linux $DIRP/
