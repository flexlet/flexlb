echo "building ..."

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/flexlb-api ../cmd/flexlb-server/main.go


#echo "copying ..."

#DIR="/root/flexlb"
#TARGET="root@<HOST>"

#scp -r run.sh /tmp/flexlb-api ../conf ../test ${TARGET}:${DIR}

#PROXY="root@<PROXY>"
#scp -r run.sh /tmp/flexlb-api ../conf ../test ${PROXY}:${DIR}
#ssh $PROXY "scp -r $DIR/* $TARGET:$DIR"
