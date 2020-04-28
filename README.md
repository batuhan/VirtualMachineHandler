# Virtual Machine Handler

This project is a wrapper around VMWare's GOVC command line tool to easily provision virtual machines & manage them via
a simple REST API. Uses `cloud-init` to provision the servers. You must have templates set up before running VMH.

Templates that we use can be found under `packer` directory.

We have a set of Packer-based template generators for various versions of Ubuntu, CentOS & Debian.

**WORK IN PROGRESS**

## Features

- Uses `cloud-init` to provision servers, providing compatibility between different linux distros.

- Backend agnostic. VMH provides its own minimal REST API, so you can attach it to your existing infrastructure easily.

- Security first. VMH generates passwords randomly and requires users to change their password on first login.

- Async by default. Non-blocking calls with HTTP callbacks. You can define auth headers and tokens for verification too.

- Archive instead of delete. VMH archives virtual machines instead of permanently deleting them, so you can keep your business data safe.

- Manage multiple datacenter configurations with a single deployment and without a database.

- Default configuration that can be shared between multiple datacenter configurations.

- Handle machine state via an easy REST API. Available actions as on, off, suspend, shutdown and reboot.

## Configuration

First need to pass vCenter configuration in your environment:

```shell script
IDENTIFIER_GOVC_INSECURE=1

IDENTIFIER_GOVC_URL=xxx
IDENTIFIER_GOVC_USERNAME=xxx
IDENTIFIER_GOVC_PASSWORD=xxx
IDENTIFIER_GOVC_DATACENTER=xxx
IDENTIFIER_GOVC_DATASTORE=xxx
IDENTIFIER_GOVC_RESOURCE_POOL=xxx

IDENTIFIER_DELETE_DIRECTORY=xxx
IDENTIFIER_TARGET_DIRECTORY=xxx
IDENTIFIER_GATEWAY=xxx
IDENTIFIER_NAMESERVERS=8.8.8.8,8.8.4.4

IDENTIFIER_WEBHOOK_URL=xxx
IDENTIFIER_WEBHOOK_AUTH_HEADER=xxx
IDENTIFIER_WEBHOOK_AUTH_TOKEN=xxx

LOCATION_IDS=IDENTIFIER1,IDENTIFIER2

HTTP_PORT=8080
POWER_OFF_TIMEOUT=1m
```

Remember to replace `IDENTIFIER` with a location ID like `AMS1` (or anything you like).

You can define multiple datacenters.
Each configuration extends the `DEFAULT` config.

If you are using a single location, you can define the defaults and use `DEFAULT` as your location ID.

You also need to set `LOCATION_IDS` with a comma separated list of location identifiers.
`DEFAULT` is enabled by default.

### Example

```shell script
DEFAULT_GOVC_INSECURE=1
AMS1_GOVC_URL=AMS1_URL
AMS2_GOVC_URL=AMS2_URL
LOCATION_IDS=AMS1,AMS2
```

In this example, both configurations will have a value of `1` for `GOVC_INSECURE`.

## Usage

VMH routes are basically RPC routes. Every action sent as a `POST` request with a JSON body. Endpoints are called actions.

For each request, you'll need to send `locationId` & `targetName`.

Calls are async by default. VMH will return only return a UUID as response.
You can use this UUID to track the progress of request with callbacks.

### env

```json
{
  "locationId": "IDENTIFIER"
}
```

### create / recreate

```json
{
  "locationId": "IDENTIFIER",
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
  "locationId": "IDENTIFIER",
  "targetName": "UbuntuTarget"
}
```

### update

```json
{
  "locationId": "CENTER2",
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
  "locationId": "IDENTIFIER",
  "targetName": "UbuntuTarget",
  "action": "shutdown"
}
```

`action` must be one of the following values: `on | off | suspend | shutdown | reboot`

#### Notes

Template names must contain one of the following values, lower/upper case doesn't matter

```
centos-7
centos-8
ubuntu
debian
```

## Creating VMware templates

The only user with login access should be `root`

`cloud-init clean` should be run after every change to template

### Ubuntu

#### Install

```
apt install python3-pip
curl -sSL https://raw.githubusercontent.com/vmware/cloud-init-vmware-guestinfo/master/install.sh | sh -
cloud-init clean
```

### Centos

#### Install

```
yum install python3
curl -sSL https://raw.githubusercontent.com/vmware/cloud-init-vmware-guestinfo/master/install.sh | sh -
```
