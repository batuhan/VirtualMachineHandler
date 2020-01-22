## env variables
```
GOVC_INSECURE=1

IDENTIFIER_GOVC_URL=xxx
IDENTIFIER_GOVC_USERNAME=xxx
IDENTIFIER_GOVC_PASSWORD=xxx
IDENTIFIER_GOVC_DATACENTER=xxx
IDENTIFIER_GOVC_DATASTORE=xxx
IDENTIFIER_GOVC_RESOURCE_POOL=xxx

IDENTIFIER_TARGET_DIRECTORY=xxx
IDENTIFIER_GATEWAY=xxx
IDENTIFIER_NAMESERVERS=8.8.8.8,8.8.4.4

WEBHOOK_URL=xxx
```
pass identifier in the request body, like below

## example requests
### env
```json
{
  "identifier": "CENTER2"
}
```
### create / recreate
```json
{
  "identifier": "IDENTIFIER",
  "template": "Ubuntu1804",
  "targetName": "UbuntuTarget",
  "cpu": 1,
  "memory": 1024,
  "diskSize": "100G",
  "sshKey": "ssh-key",
  "ipToAssign": "1.1.1.1"
}
```
### delete
```json
{
  "identifier": "IDENTIFIER",
  "targetName": "UbuntuTarget"
}
```
### update
```json
{
  "identifier": "CENTER2",
  "targetName": "UbuntuTarget",
  "cpu": 1,
  "memory": 1024,
  "diskSize": "50G"
}
```
`cpu`, `memory` and `diskSize` fields can be omitted, only provided values will be updated 
### state
```json
{
  "identifier": "IDENTIFIER",
  "targetName": "UbuntuTarget",
  "action": "shutdown"
}
```
`action` must be one of the following values: `on | off | suspend | shutdown | reboot`
#### other
template names must contain one of the following values, lower/upper case doesn't matter
```
centos-7
centos-8
ubuntu
debian
```

## while creating vmware templates

the only user with login access should be `root`

`cloud-init clean` should be run after every change to template

### Ubuntu

#### to install
```
apt install python3-pip
curl -sSL https://raw.githubusercontent.com/vmware/cloud-init-vmware-guestinfo/master/install.sh | sh -
cloud-init clean
```
### Centos

#### to install
```
yum install python3
curl -sSL https://raw.githubusercontent.com/vmware/cloud-init-vmware-guestinfo/master/install.sh | sh -
```
`cloud-init clean` after every change to template
