# FlexLB API Server

Flexible load balancer API server to control keepalived and haproxy

## Build

### Clone code

```sh
git clone https://github.com/flexlet/flexlb.git
```

### Build binary

#### For Linux
```sh
cd flexlb

# edit build/profile to config options

# generate certs
sh build/gencrt.sh

# build tarball
sh build/build.sh
```

## Run

### Install and start Keepalived

#### For EulerOS 2.9
```sh
rpm -ivh net-snmp-5.8-8.h6.eulerosv2r9.x86_64.rpm  net-snmp-libs-5.8-8.h6.eulerosv2r9.x86_64.rpm
rpm -ivh keepalived-2.0.20-16.h5.eulerosv2r9.x86_64.rpm
systemctl enable keepalived
systemctl start keepalived
```

### Install HAProxy

#### For EulerOS 2.9
```sh
rpm -ivh haproxy-2.0.14-1.eulerosv2r9.x86_64.rpm
```

### Run FlexLB API server
```sh
# copy and extract tarbal
mkdir flexlb-<version> && cd flexlb-<version>
tar -zxf ../flexlb-<version>.tar.gz

# check install options
sh install.sh -h

# install
sh install.sh [options]
```

## Test

### Prepare backend servers

For example, create 3 nginx server:
```
192.168.1.141:30080
192.168.1.142:30080
192.168.1.143:30080
```

### Prepare instance config template

Prepare instance config for API post, for example: 
```
hack/instance_template.json
```

### Test instance create, list, modify, get, start, stop, delete

```sh
# get ready status
sh get_status.sh

# create instances
sh create_instance.sh inst1 192.168.2.1
sh create_instance.sh inst2 192.168.2.2

# list instances
sh list_instances.sh
sh list_instances.sh inst1

# modify instance
# edit instance_template.json to add or remove backend servers
sh modify_instance.sh inst1 192.168.2.3

# show instance changes
sh get_instance.sh inst1

# stop & start instance
sh stop_instance.sh inst1
sh start_instance.sh inst1

# delete instances
sh delete_instance.sh inst1
sh delete_instance.sh inst2
```
