# Virtual Machine Handler

This project is a wrapper around VMWare's GOVC command line tool to easily provision virtual machines & manage them via
a simple REST API. Uses `cloud-init` to provision the servers. You must have templates set up before running VMH. 

We have a set of Packer-based template generators for various versions of Ubuntu, CentOS & Debian.

**WORK IN PROGRESS**

## Configuration

First need to pass vCenter configuration in your environment:
 
```
IDENTIFIER_GOVC_INSECURE=1

IDENTIFIER_GOVC_URL=xxx
IDENTIFIER_GOVC_USERNAME=xxx
IDENTIFIER_GOVC_PASSWORD=xxx
IDENTIFIER_GOVC_DATACENTER=xxx
IDENTIFIER_GOVC_DATASTORE=xxx
IDENTIFIER_GOVC_RESOURCE_POOL=xxx

IDENTIFIER_TARGET_DIRECTORY=xxx
IDENTIFIER_GATEWAY=xxx
IDENTIFIER_NAMESERVERS=8.8.8.8,8.8.4.4

IDENTIFIER_WEBHOOK_URL=xxx
IDENTIFIER_WEBHOOK_AUTH_HEADER=xxx
IDENTIFIER_WEBHOOK_AUTH_TOKEN=xxx
```

Remember to replace `IDENTIFIER` with a location ID like `AMS1` (or anything you like). 
You can also define a `DEFAULT` location.

If you are using a single location, you can define the defaults and use `DEFAULT` as your location ID.

You also need to set `LOCATION_IDS` with a comma separated list of location identifiers

You can set the port the server runs on with `HTTP_PORT`.

## Usage

VMH routes are basically RPC routes. Every action sent as a `POST` request with a JSON body. Endpoints are called actions.

For each request, you'll need to send `LocationId` & `TargetName`. 

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
