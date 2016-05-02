#!/usr/bin/env bash

#sudo apt-get update -y
#sudo apt-get install -y libc6 libstdc++6
#sudo iptables -A INPUT -p tcp -m multiport --dports 11211,11210,11209,4369,8091,8092,8093,9998,18091,18092,11214,11215  -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 11211 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 11210 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 11209 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 4369 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8091 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8092 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8093 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 9998 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 18091 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 18092 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 11214 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 11215 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 9100:9105 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 21100:21299 -j ACCEPT
sudo iptables -L
sudo wget http://packages.couchbase.com/releases/4.5.0-DP1/couchbase-server-enterprise_4.5.0-DP1-ubuntu14.04_amd64.deb
sudo dpkg -i couchbase-server-enterprise_4.5.0-DP1-ubuntu14.04_amd64.deb
sudo service couchbase-server restart
sleep 15
/opt/couchbase/bin/couchbase-cli cluster-init -c $HOSTNAME:8091 --cluster-username=Administrator --cluster-password=password --cluster-ramsize=512
/opt/couchbase/bin/couchbase-cli bucket-create -c $HOSTNAME:8091 -u Administrator -p password --bucket=documents --bucket-type=couchbase --bucket-ramsize=512 --wait
