#!/usr/bin/env bash

if [ ! -f downloads/couchbase_server-4.5.0-DP1.deb ]; then
    wget http://packages.couchbase.com/releases/4.5.0-DP1/couchbase-server-enterprise_4.5.0-DP1-ubuntu14.04_amd64.deb -O downloads/couchbase_server-4.5.0-DP1.deb
fi

dpkg-deb -x downloads/couchbase_server-4.5.0-DP1.deb $HOME
cd $HOME/opt/couchbase
./bin/install/reloc.sh `pwd`
./bin/couchbase-server -- -noinput -detached
#sudo service couchbase-server restart
sleep 20
./bin/couchbase-cli cluster-init -c localhost:8091 --cluster-username=Administrator --cluster-password=password --cluster-ramsize=512
./bin/couchbase-cli bucket-create -c localhost:8091 -u Administrator -p password --bucket=testing --bucket-type=couchbase --bucket-ramsize=512 --wait
