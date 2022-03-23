CONF=/etc/flexlb
export HOST=localhost
export PORT=8080
export TLS_HOST=localhost
export TLS_PORT=8443
export TLS_CERTIFICATE=${CONF}/certs/server.crt
export TLS_PRIVATE_KEY=${CONF}/certs/server.key
export TLS_CA_CERTIFICATE=${CONF}/certs/ca.crt

chmod +x flexlb-api

./flexlb-api --config-file=conf/flexlb-api-config.yaml --debug