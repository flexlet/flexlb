if [ $# -ne 3 ]; then
  echo "$0 <TEMPLATE> <NAME> <VIP>"
  exit 1
fi
TEMPLATE=$1
NAME=$2
VIP=$3

URI="http://localhost:8080/flexlb/v1/instances"

DATA=$(sed "s/<NAME>/${NAME}/g; s/<VIP>/${VIP}/g" ${TEMPLATE})

curl -H "Content-Type: application/json" -X POST -d "$DATA" $URI
