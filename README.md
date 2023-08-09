# cubefs-dashboard

# Introduction
cubefs Console (cfs-gui) is a graphical management tool for cubeFS system, which is designed to meet the needs of system administrators and ordinary users for daily file management. It is an independent web service that is more intuitive and easier to get started with than cfs-cli.

# Architecture
The cfs-gui backend interacts directly with the CubeFS cluster, manages multiple clusters simultaneously, and stores its own data in the database. The permission module can assign accounts to different application scenarios. cfs-gui has four main functional modules, as follows:<br>
- Logical resource management: volume, data-partition, meta-partition, erasure-coding stripe, etc
- Physical resource management: nodes, disks, etc
- Authorization management: users and roles
- File management: file upload, download, list

<img src="https://github.com/cubefs/cubefs-dashboard/blob/main/pictures/architecture.png" align=center/>

# Features
- Multi-cluster management: the dashboard can be deployed independently from the storage cluster and can manage multiple storage clusters<br>
- Flexible authorization management: Each user can be associated with multiple roles, and each role has access permissions to a set of resources
- Volume creation and authorization: it supports authorizing other tenants to read and write the volumes created by oneself
- Resource migration: it supports resource decommissioning or migration at the granularity of partition, disk, or machine
- File management: List, upload, download<br>

# Compilation
## Requirements:
node version: v12-v14<br>
npm version: the npm in node package<br>
go version: v1.16 and later<br>

```
## clone code
git clone https://github.com/cubefs/cubefs-dashboard.git

## build
make

## product
ll bin/
total 26256
-rwxr-xr-x 1 root root 26876984 Jul 27 17:26 cfs-gui
-rw-r--r-- 1 root root     1009 Jul 27 17:26 config.yml
drwxr-xr-x 6 root root     4096 Jul 27 17:27 dist
```
# Preview
## Cluster overview
<img src="https://github.com/cubefs/cubefs-dashboard/blob/main/pictures/cluster_overview.png" align=center/>

## Role management
<img src="https://github.com/cubefs/cubefs-dashboard/blob/main/pictures/role_management.png" align=center/>

## File management
<img src="https://github.com/cubefs/cubefs-dashboard/blob/main/pictures/file_management.png" align=center/>

