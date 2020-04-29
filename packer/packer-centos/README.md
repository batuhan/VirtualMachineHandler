# CentOS packer templates

### How to generate ks\*.cfg files

```shell script
./gen-ks.sh 8
```

Change 8 with the version you want to install. Supported versions: 7, 8

### Vagrant Virtualbox build example

    packer build -only=virtualbox-iso --var 'build_version=1.0.0' \
    	-var-file=centos-8-vagrant.json centos-vagrant.json

Vagrant boxes will be automaticly uploaded to [Vagrant Cloud](https://app.vagrantup.com/). This requires the `VAGRANT_CLOUD_TOKEN` environment variable to be set with a valid [authentication token](https://app.vagrantup.com/settings/security), otherwise the upload will fail.

### VMware vSphere [Jetbrains packer-builder-vsphere](https://github.com/jetbrains-infra/packer-builder-vsphere) build example

    cp vsphere_environment.sh.dist vsphere_environment.sh
    edit vsphere_environment.sh
    source vsphere_environment.sh
    packer build --var 'whiteout=false' -var-file=centos-8.json centos-vsphere.json

#### Required environment variables

    PACKER_VSPHERE_VCENTER_SERVER='hostname_or_IP'
    PACKER_VSPHERE_ESXI_HOST='hostname_or_IP'
    PACKER_VSPHERE_DATACENTER='datacenter'
    PACKER_VSPHERE_RESOURCE_POOL='resource_pool'
    PACKER_VSPHERE_USERNAME='username'
    PACKER_VSPHERE_PASSWORD='password'
    PACKER_VSPHERE_DATASTORE='datastore_name'
    PACKER_VSPHERE_DATASTORE_ISO='datastore with ISO directory'
    PACKER_VSPHERE_NETWORK='network_name'
    PACKER_VSPHERE_VM_VERSION='13'

### VMware vSphere native packer build example

    cp vsphere_environment.sh.dist vsphere_environment.sh
    edit vsphere_environment.sh
    source vsphere_environment.sh
    packer build -var-file=centos-8.json centos-vsphere-native.json

#### Required environment variables

    PACKER_VMWARE_HOST='hostname_or_IP'
    PACKER_VMWARE_USERNAME='username'
    PACKER_VMWARE_PASSWORD='password'
    PACKER_VMWARE_DATASTORE='datastore_name'
    PACKER_VMWARE_NETWORK='network_name'
    PACKER_VMWARE_MAC='00:cc:aa:bb:ee:ee'
    PACKER_VMWARE_VM_VERSION='13'
