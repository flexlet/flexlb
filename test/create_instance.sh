if [ $# -ne 2 ]; then
  echo "$0 <NAME> <VIP>"
  exit 1
fi
NAME=$1
VIP=$2

URI="http://localhost:8080/flexlb/v1/instances"

TEMPLATE="instance_template.json"
DATA=$(sed "s/<NAME>/${NAME}/g; s/<VIP>/${VIP}/g" ${TEMPLATE})

curl -H "Content-Type: application/json" -X POST -d "$DATA" $URI
