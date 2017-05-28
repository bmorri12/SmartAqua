#!/bin/sh

SERVICES='registry devicemanager mqttaccess controller httpaccess apiprovider'

echo "stopping all services..."
docker stop `echo $SERVICES`
docker rm `echo $SERVICES`

sudo docker run --link mysql mysql sh -c 'exec mysql -hmysql -uroot -e"CREATE DATABASE SeaWater"'

echo "starting registry..."
docker run -d --name registry --link etcd --link mysql xfocus/smartaquaculture registry -etcd http://etcd:2379 -rpchost internal:20034 -aeskey ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP -dbhost mysql -dbname SeaWater -dbport 3306 -dbuser root -loglevel debug

echo "starting devicemanager..."
docker run -d --name devicemanager --link etcd --link redis xfocus/smartaquaculture devicemanager -etcd http://etcd:2379 -rpchost internal:20033 -redishost redis:6379 -loglevel debug

echo "starting controller..."
docker run -d --name controller --link etcd --link mongo xfocus/smartaquaculture controller -etcd http://etcd:2379 -mongohost mongo -rpchost internal:20032 -loglevel debug

echo "starting mqttaccess..."
docker run -d --name mqttaccess -p 1883:1883 --link etcd xfocus/smartaquaculture mqttaccess -etcd http://etcd:2379 -tcphost :1883 -usetls -keyfile /go/src/github.com/bmorri12/SmartAqua/pkg/server/testdata/key.pem -cafile /go/src/github.com/bmorri12/SmartAqua/pkg/server/testdata/cert.pem -rpchost internal:20030 -loglevel debug

echo "starting httpaccess..."
docker run -d --name httpaccess -p 443:443 --link etcd --link redis xfocus/smartaquaculture httpaccess -etcd http://etcd:2379 -httphost :443 -redishost redis:6379 -usehttps -keyfile /go/src/github.com/bmorri12/SmartAqua/pkg/server/testdata/key.pem -cafile /go/src/github.com/bmorri12/SmartAqua/pkg/server/testdata/cert.pem -loglevel debug

echo "starting apiprovider..."
docker run -d --name apiprovider -p 8888:8888 --link etcd xfocus/smartaquaculture apiprovider -etcd http://etcd:2379 -httphost :8888 -loglevel debug

sudo docker run -it --name pdcfg -v `echo $(pwd)`:/root --link etcd bmorri12/SmartAqua pdcfg -etcd http://etcd:2379