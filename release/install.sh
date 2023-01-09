#!/bin/bash

# default options
DEFAULT_NODE_NAME=${HOSTNAME}
DEFAULT_HOST=127.0.0.1
DEFAULT_PORT=8080
DEFAULT_TLS_HOST=$(hostname -i)
DEFAULT_TLS_PORT=8443
DEFAULT_CLUSTER_NAME="FLEXLB"
DEFAULT_ADVERTIZE_PORT=8001
DEFAULT_CLUSTER_SECRET_KEY="FlexLBDefaultKey"

# installation target
BIN="/usr/local/bin"
ETC="/etc/flexlb"
SVC="/usr/lib/systemd/system/flexlb.service"

function fn_print_help() {
    echo "$(basename $0) [options]
    Options:
        -N NAME                optional, node name, default: ${DEFAULT_NODE_NAME}
        -H HOST                optional, flexlb http listening host, default: ${DEFAULT_HOST}
        -P PORT                optional, flexlb http listening port, default: ${DEFAULT_PORT}
        -h TLS_HOST            optional, tls listening host, default: ${DEFAULT_TLS_HOST}
        -p TLS_PORT            optional, tls listening port, default: ${DEFAULT_TLS_PORT}
        -n CLUSTER_NAME        optional, flexlb cluster name, default: ${DEFAULT_CLUSTER_NAME}
        -e CLUSTER_ENDPOINT    required, flexlb cluster public endpoint (<cluster_vip>:<port>)
        -a CLUSTER_ADVERTIZE   optional, flexlb cluster advertize endpoint (<node_private_ip>:<port>), default: ${DEFAULT_TLS_HOST}:${DEFAULT_ADVERTIZE_PORT}
        -m CLUSTER_MEMBERS     optional, flexlb cluster members, comma seperate advertize endpoints (<node1_private_ip>:<port>,<node1_private_ip>:<port>), default local: CLUSTER_ADVERTIZE
        -s CLUSTER_SECRET_KEY  optional, flexlb cluster secret key, default: ${DEFAULT_CLUSTER_SECRET_KEY}
    "
    exit -1
}

function fn_validate_params() {
    while getopts N:H:P:h:p:n:e:a:m:s: flag
    do
        case "${flag}" in
            N) NAME=${OPTARG};;
            H) HOST=${OPTARG};;
            P) PORT=${OPTARG};;
            h) TLS_HOST=${OPTARG};;
            p) TLS_PORT=${OPTARG};;
            n) CLUSTER_NAME=${OPTARG};;
            e) CLUSTER_ENDPOINT=${OPTARG};;
            a) CLUSTER_ADVERTIZE=${OPTARG};;
            m) CLUSTER_MEMBERS=${OPTARG};;
            s) CLUSTER_SECRET_KEY=${OPTARG};;
            ?) fn_print_help
        esac
    done

    if [ "${NAME}" == "" ]; then NAME=${DEFAULT_NODE_NAME}; fi
    if [ "${HOST}" == "" ]; then HOST=${DEFAULT_HOST}; fi
    if [ "${PORT}" == "" ]; then PORT=${DEFAULT_PORT}; fi
    if [ "${TLS_HOST}" == "" ]; then TLS_HOST=${DEFAULT_TLS_HOST}; fi
    if [ "${TLS_PORT}" == "" ]; then TLS_PORT=${DEFAULT_TLS_PORT}; fi
    if [ "${CLUSTER_NAME}" == "" ]; then CLUSTER_NAME=${DEFAULT_CLUSTER_NAME}; fi
    if [ "${CLUSTER_ENDPOINT}" == "" ]; then fn_print_help; fi
    if [ "${CLUSTER_ADVERTIZE}" == "" ]; then CLUSTER_ADVERTIZE="${DEFAULT_TLS_HOST}:${DEFAULT_ADVERTIZE_PORT}"; fi
    if [ "${CLUSTER_MEMBERS}" == "" ]; then CLUSTER_MEMBERS=${CLUSTER_ADVERTIZE}; fi
    if [ "${CLUSTER_SECRET_KEY}" == "" ]; then CLUSTER_SECRET_KEY=${DEFAULT_CLUSTER_SECRET_KEY}; fi
}

function fn_main() {
    fn_validate_params $@

    # stop old service
    active=$(systemctl is-active flexlb)
    if [ "${active}" == "active" ]; then
        systemctl stop flexlb
        systemctl disable flexlb
        rm ${SVC}
        systemctl daemon-reload
    fi

    # clean old config
    if [ -d "${ETC}" ]; then rm -rf ${ETC}; fi
    
    # generate config
    mkdir -p ${ETC}
    cp -r etc/certs ${ETC}/
    cp etc/service.conf ${ETC}/
    sed -e "s/^name:.*$/name: ${NAME}/" \
        -e "s/^host:.*$/host: ${HOST}/" \
        -e "s/^port:.*$/port: ${PORT}/" \
        -e "s/^tls_host:.*$/tls_host: ${TLS_HOST}/" \
        -e "s/^tls_port:.*$/tls_port: ${TLS_PORT}/" \
        -e "s/  name:.*$/  name: ${CLUSTER_NAME}/" \
        -e "s/  endpoint:.*$/  endpoint: \"${CLUSTER_ENDPOINT}\"/" \
        -e "s/  advertize:.*$/  advertize: \"${CLUSTER_ADVERTIZE}\"/" \
        -e "s/  member:.*$/  member: \"${CLUSTER_MEMBERS}\"/" \
        -e "s/  secret_key:.*$/  secret_key: \"${CLUSTER_SECRET_KEY}\"/" \
        etc/config.yaml > ${ETC}/config.yaml

    # install binary
    cp bin/flexlb  ${BIN}/
    chmod +x ${BIN}/flexlb

    # install service
    cp systemd/flexlb.service ${SVC}
    systemctl daemon-reload
    systemctl enable flexlb
    systemctl start flexlb
    sleep 0.5
    systemctl status flexlb
}

fn_main $@

