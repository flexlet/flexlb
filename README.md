# FlexLB API Server

Flexible load balancer API server to control keepalived and haproxy

## Build

### Generate code

```sh
swagger generate server -f swagger/flexlb-api-spec.yaml
```

### Build binary

#### For Linux
```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/flexlb-api cmd/flexlb-server/main.go
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

### Generage self-signed certificate

#### Generate CA key and CA certs

```sh
mkdir -p  /etc/flexlb/certs/
cd /etc/flexlb/certs/

DNS_NAME="example.com"

openssl genrsa -out ca.key 2048
openssl req -new -out ca.csr -key ca.key -subj "/CN=${DNS_NAME}"
openssl x509 -req -in ca.csr -out ca.crt -signkey ca.key -CAcreateserial -days 3650
```
#### Generate server key and certs

```sh
openssl genrsa -out server.key 2048
openssl req -new -out server.csr -key server.key -subj "/CN=${DNS_NAME}"
openssl x509 -req -in server.csr -out server.crt -signkey server.key -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650
```

### Run FlexLB API server
```sh
CONF=/etc/flexlb
export HOST=localhost
export PORT=8080
export TLS_HOST=localhost
export TLS_PORT=8443
export TLS_CERTIFICATE=${CONF}/certs/server.crt
export TLS_PRIVATE_KEY=${CONF}/certs/server.key
export TLS_CA_CERTIFICATE=${CONF}/certs/ca.crt

# copy flexlb-api and conf/flexlb-api-config.yaml here
# run flexlb-api server in debug mode
chmod +x flexlb-api
./flexlb-api --config-file=conf/flexlb-api-config.yaml --debug
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
test/instance_template.json
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