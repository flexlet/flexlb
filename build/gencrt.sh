#!/bin/bash

BASE_DIR="$(readlink -f $(dirname $0)/..)"
BUILD_DIR="${BASE_DIR}/build"
RELEASE_DIR="${BASE_DIR}/release"

# inject build options
source ${BUILD_DIR}/profile

# build target
TARGET_DIR="${BUILD_DIR}/target"
CERTS_DIR="${TARGET_DIR}/certs"

# clean certs
rm -rf ${CERTS_DIR}

# create certs directories
mkdir -p ${CERTS_DIR}

cd ${CERTS_DIR}

# Generate ca key and certs
openssl genrsa -out ca.key 2048
openssl req -new -out ca.csr -key ca.key -subj "${CERT_ISSUER}"
openssl x509 -req -in ca.csr -out ca.crt -signkey ca.key -CAcreateserial -days ${CERT_EXPIRE}

# Generate server cert config
cat <<EOF > server.cnf
[req]
distinguished_name = req_distinguished_name
req_extensions = req_ext
prompt = no

[req_distinguished_name]
CN  = ${SERVER_DNS_NAME}

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${SERVER_DNS_NAME}
EOF

# Generate server key and certs
openssl genrsa -out server.key 2048
openssl req -new -out server.csr -key server.key -config server.cnf
openssl x509 -req -days ${CERT_EXPIRE} -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extensions req_ext -extfile server.cnf

# Generate client cert config
cat <<EOF > client.cnf
[req]
distinguished_name = req_distinguished_name
req_extensions = req_ext
prompt = no

[req_distinguished_name]
CN  = ${CLIENT_DNS_NAME}

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${CLIENT_DNS_NAME}
EOF

#### Generate client key and certs
openssl genrsa -out client.key 2048
openssl req -new -out client.csr -key client.key -config client.cnf
openssl x509 -req -days ${CERT_EXPIRE} -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -extensions req_ext -extfile client.cnf

# copy to release dir
cp ca.crt server.crt server.key ${RELEASE_DIR}/etc/certs/

cd ..

