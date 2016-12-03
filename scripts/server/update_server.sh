#!/usr/bin/env bash

printf "%-60s" "Syncing dependencies"
govendor sync
if [ $? -eq 0 ]; then
    printf "%-10s\n" "OK"
else
    printf "%-10s\n" "FAIL"
fi

printf "%-60s" "Stopping CodeCollaborate daemon"
sudo systemctl stop CodeCollaborate.service
if [ $? -eq 0 ]; then
    printf "%-10s\n" "OK"
else
    printf "%-10s\n" "FAIL"
fi

printf "%-60s" "Building server binary"
go build -o CodeCollaborateServer
if [ $? -eq 0 ]; then
    printf "%-10s\n" "OK"
else
    printf "%-10s\n" "FAIL"
fi

printf "%-60s\n" "Replacing server binary and setting permissions"
sudo mv CodeCollaborateServer /CodeCollaborate/Server
if [ $? -eq 0 ]; then
    printf "\t%-10s\n" "COPY OK"
else
    printf "\t%-10s\n" "COPY FAIL"
fi
sudo chmod +x /CodeCollaborate/Server
if [ $? -eq 0 ]; then
    printf "\t%-10s\n" "CHMOD OK"
else
    printf "\t%-10s\n" "CHMOD FAIL"
fi

printf "%-60s" "Starting service"
sudo systemctl start CodeCollaborate.service
if [ $? -eq 0 ]; then
    printf "%-10s\n" "OK"
else
    printf "%-10s\n" "FAIL"
fi
