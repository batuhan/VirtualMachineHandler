#!/bin/bash -eux

source /etc/os-release
# Install Python.
yum -y install python3-pip
alternatives --set python /usr/bin/python3

# Install Ansible.
pip3 install ansible
yum -y install git

# Install cloud-init
if [[ $VERSION_ID == 7 ]]; then
  yum -y install epel-release
  yum -y install python-pip
fi
yum -y install python3 python3-pip cloud-init cloud-utils-growpart
curl -sSL https://raw.githubusercontent.com/vmware/cloud-init-vmware-guestinfo/master/install.sh | sh -
cloud-init clean
