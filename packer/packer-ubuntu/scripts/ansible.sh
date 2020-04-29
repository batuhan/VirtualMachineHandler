#!/bin/bash -eux

export DEBCONF_NONINTERACTIVE_SEEN=true
export DEBIAN_FRONTEND=noninteractive
export UCF_FORCE_CONFOLD=1

# Install Ansible repository.
apt -y update
apt -y install software-properties-common git curl python3-pip cloud-init cloud-initramfs-growroot
apt-add-repository ppa:ansible/ansible

# Install Ansible and other dependencies

# This
apt-get -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" -y -qq install ansible


# Install VMWare integration to cloud-init
curl -sSL https://raw.githubusercontent.com/vmware/cloud-init-vmware-guestinfo/master/install.sh | sh -
cloud-init clean