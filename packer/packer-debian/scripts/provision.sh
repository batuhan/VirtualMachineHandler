#!/bin/bash -eux

# Add PPA
echo "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" >/etc/apt/sources.list.d/ansible.list

# Required by Debian 9; otherwise apt-key will fail.
apt-get -y install dirmngr

apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367

# apt cdrom after-install fix
sed -i '/^deb cdrom:/s/^/#/' /etc/apt/sources.list

# Install Ansible & other dependencies.
apt-get -y update
apt-get -y install ansible git curl python3-pip python3-jsonschema wget

# Install cloud-init
#curl http://archive.ubuntu.com/ubuntu/pool/main/c/cloud-init/cloud-init_19.3-41-gc4735dd3-0ubuntu1~18.04.1_all.deb >/tmp/cloud-init.deb
#curl http://archive.ubuntu.com/ubuntu/pool/main/c/cloud-initramfs-tools/cloud-initramfs-growroot_0.25ubuntu1.12.04.1_all.deb >/tmp/cloud-initramfs-growroot.deb
#dpkg -i /tmp/cloud-init.deb >/dev/null 2>&1
#dpkg -i /tmp/cloud-initramfs-growroot.deb >/dev/null 2>&1
#apt-get --fix-broken -y install
apt-get -y install cloud-init cloud-initramfs-growroot python3 python3-pip python-dev python-pip

# Install VMWare integration to cloud-init
curl -sSL https://raw.githubusercontent.com/vmware/cloud-init-vmware-guestinfo/master/install.sh | sh -

#cloud-init clean
