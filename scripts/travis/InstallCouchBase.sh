#!/usr/bin/env bash

if [ ! -f ~/downloads/couchbase_server-4.5.0-DP1.deb ]; then
    wget http://packages.couchbase.com/releases/4.5.0-DP1/couchbase-server-enterprise_4.5.0-DP1-ubuntu14.04_amd64.deb -O ~/downloads/couchbase_server-4.5.0-DP1.deb
fi

dpkg-deb -x ~/downloads/couchbase_server-4.5.0-DP1.deb $HOME
cd $HOME/opt/couchbase
./bin/install/reloc.sh `pwd`
./bin/couchbase-server -- -noinput -detached

sleep 20
./bin/couchbase-cli cluster-init -c localhost:8091 --cluster-username=username --cluster-password=password --cluster-ramsize=512
./bin/couchbase-cli bucket-create -c localhost:8091 -u username -p password --bucket-password=password --bucket=cc --bucket-type=couchbase --bucket-ramsize=412 --wait
./bin/couchbase-cli bucket-create -c localhost:8091 -u username -p password --bucket-password=password --bucket=cc_scrunching_locks --bucket-type=couchbase --bucket-ramsize=100 --enable-index-replica=0 --bucket-replica=0 --enable-flush=1 --wait
