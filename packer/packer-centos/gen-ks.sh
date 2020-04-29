#!/usr/bin/env sh

FILE=./http/ks${1}.cfg
rm "${FILE}" > /dev/null 2>&1
cp ./http/template-ks"${1}".cfg "${FILE}"
sed -i "s/NETWORK_IP/${PACKER_VSPHERE_NETWORK_IP}/" "${FILE}"
sed -i "s/NETWORK_GATEWAY/${PACKER_VSPHERE_NETWORK_GATEWAY}/" "${FILE}"
sed -i "s/NETWORK_MASK/${PACKER_VSPHERE_NETWORK_MASK}/" "${FILE}"
sed -i "s/NETWORK_NAMESERVER1/${DEFAULT_DNS1}/" "${FILE}"
sed -i "s/NETWORK_NAMESERVER2/${DEFAULT_DNS2}/" "${FILE}"
sed -i "s/NETWORK_NAMESERVER/${DEFAULT_DNS1},${DEFAULT_DNS2}/" "${FILE}"
