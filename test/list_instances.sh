URI="http://localhost:8080/flexlb/v1/instances"

if [ $# -eq 1 ]; then
  URI="${URI}?name=$1"
fi

curl $URI
