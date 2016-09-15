#!/usr/bin/env bash

# Setup MySQL DB & Tables
printf "\n"
printf "Setting up MySQL DB & Tables\n"
docker exec -i CodeCollaborate_MySQL mysql --protocol=tcp -uroot -pCodeCollaborate1234 < ../config/defaults/mysqlSetup.sql
docker exec -i CodeCollaborate_MySQL mysql --protocol=tcp -uroot -pCodeCollaborate1234 < ../config/defaults/mysqlSchemaSetup.sql