#!/bin/bash

# Requires jq, wget & govc
# Remember to run `source vsphere_environment.sh` first

govc datastore.mkdir ISO || true

 
IFS='-'
read -ra OS_KEY_PARTS <<< "$OS_KEY"
IFS=' '

ISO_DATA=`cat packer-$OS_KEY_PARTS/$OS_KEY.json`
ISO_URL=`echo $ISO_DATA | jq -r .iso_download_url`
ISO_FILE=`echo $ISO_DATA | jq -r .iso_name`
ISO_CHECKSUM=`echo $ISO_DATA | jq -r .iso_checksum`

[ -f $ISO_FILE ] && echo "$ISO_FILE exist" || wget $ISO_URL

# @TODO: Check checksum
govc datastore.upload $ISO_FILE "ISO/$ISO_FILE"