# swarm-volume
Some Docker Volume Plugins that provision really simple swarm aware storage.

Hosted on [Docker Hub](https://hub.docker.com/r/chrisbecke/swarm-volume)

## Overview

This docker volume plugin provides a simple way to use a single network share as a root storage for multiple docker volumes. To use this plugin you need already to have
created a share, currently glusterfs, nfs or AWS EFS are supported. Then, install the apprpriate plugin on your docker instance or docker swarm nodes and configure the 
details to your specific root share.

This plugin will then expose all subfolders of that share as docker volumes, and `docker volume` commands can be used to create or remove volumes as required.

## Installing

Install the plugin by assigning it an alias. Each plugin alias can be configured with different parameters either at creation, or later.

Installing the nfs variant requires NFS_DEVICE in the typical `<server>:/<path>` format:

```bash
docker plugin install chrisbecke/swarm-volume:nfs --alias swarmvol --grant-all-permissions NFS_DEVICE="10.0.0.11:/"
docker plugin set swarmvol NFS_OPTIONS=nfsvers=4.1,timeo=600,retrans=2,noresvport
```

The alias as it is important when creating volumes from docker or compose.

e.g. to create a volume on the mounted nfs share:

```bash
docker volume create --driver swarmvol my-shared-volume
```

Or from compose. 

```yaml
volumes:
  data:
    driver: swarmvol
  global:
    driver: swarmvol
    name: global

services:
```

## Configuration

The prefix indicates which plugin variant the setting applies to:

KEY | Required | Description
--- | --- | ---
NFS_TYPE | no | Defaults to `nfs4`
NFS_DEVICE | yes | `<nfs-server>:/<path>` to the nfs share.
NFS_OPTIONS | no | -o options
EFS_ID | yes | The EFS id to mount
GFS_VOLUME | yes | the name of the glusterfs volume that was created. e.g. `gfs-docker-0`
GFS_SERVERS | yes | volfile servers by dns name or ip. e.g. `10.20.0.4,10.20.0.5,10.20.0.6`

## Building

This is sadly not a seamless process. You will need to edit the Makefile appropriate to your environment and then execute `make build`. Currently only x86 is supported, other platforms will require a custom build.
