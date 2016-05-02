#!/usr/bin/env bash

sudo wget http://packages.couchbase.com/releases/4.5.0-DP1/couchbase-server-enterprise_4.5.0-DP1-ubuntu14.04_amd64.deb
sudo dpkg -i couchbase-server-enterprise_4.5.0-DP1-ubuntu14.04_amd64.deb
sudo service couchbase-server restart
sleep 5
/opt/couchbase/bin/couchbase-cli cluster-init -c localhost:8091 --cluster-username=Administrator --cluster-password=password --cluster-ramsize=512
/opt/couchbase/bin/couchbase-cli bucket-create -c localhost:8091 -u Administrator -p password --bucket=documents --bucket-type=couchbase --bucket-ramsize=512 --wait
