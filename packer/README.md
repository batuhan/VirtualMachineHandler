# Packer Templates

Generates VMware templates for centos, ubuntu and debian systems.

Details for each implementation can be found in their respective README's.

### Compatibility

Compatible with packer 1.4.5 as of writing (2020-02-18)

### Upload ISO's

ISO's can be uploaded to VMVware with the following command

```shell script
govc datastore.upload distro.iso /ISO/distro.iso
```
