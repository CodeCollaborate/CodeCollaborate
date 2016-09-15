#!/usr/bin/env bash

printf "\n"
printf "Creating data volume for MySQL Container\n"
docker volume create --name CodeCollaborate_MySQL_Data

printf "\n"
printf "Starting MySQL Container\n"
docker run --name CodeCollaborate_MySQL -v CodeCollaborate_MySQL_Data:/var/lib/mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=CodeCollaborate1234 -d mysql

printf "\n"
printf "Starting RabbitMQ Container\n"
docker run --name CodeCollaborate_RabbitMQ -p 5672:5672 -p 15672:15672 -d --hostname CodeCollaborate-RabbitMQ rabbitmq:management

printf "\n"
printf "Creating data volume for Couchbase Container\n"
docker volume create --name CodeCollaborate_Couchbase_Data

printf "\n"
printf "Starting Couchbase Container\n"
docker run -d --name CodeCollaborate_Couchbase -v CodeCollaborate_Couchbase_Data:/opt/couchbase/var/lib/couchbase/data -p 8091-8094:8091-8094 -p 11210-11300:11210-11300 couchbase
