required env variables
```
GOVC_INSECURE=1

IDENTIFIER_GOVC_URL=xxx
IDENTIFIER_GOVC_USERNAME=xxx
IDENTIFIER_GOVC_PASSWORD=xxx
IDENTIFIER_GOVC_DATACENTER=xxx
IDENTIFIER_GOVC_DATASTORE=xxx
IDENTIFIER_GOVC_RESOURCE_POOL=xxx

IDENTIFIER_TARGET_DIRECTORY=xxx
```
pass identifier in the request body, like this:
```json
{
  "identifier": "IDENTIFIER",
  "template": "Ubuntu1804",
  "targetName": "UbuntuTarget",
  "cpu": 1,
  "memory": 1024,
  "diskSize": "100G",
  "sshKey": "ssh-key"
}
```
