if [ $# -ne 1 ]; then
  echo "$0 <NAME>"
  exit 1
fi
NAME=$1

URI="http://localhost:8080/flexlb/v1/instances/${NAME}"

curl -X DELETE $URI
