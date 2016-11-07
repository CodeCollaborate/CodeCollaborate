#!/usr/bin/env bash

printf "\n"
printf "Deleting MySQL container and volume\n"
docker rm -f CodeCollaborate_MySQL
docker volume rm CodeCollaborate_MySQL_Data

printf "\n"
printf "Deleting RabbitMQ container and volume\n"
docker rm -f CodeCollaborate_RabbitMQ

printf "\n"
printf "Deleting Couchbase container and volume\n"
docker rm -f CodeCollaborate_Couchbase
docker volume rm CodeCollaborate_Couchbase_Data