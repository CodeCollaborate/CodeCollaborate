#!/usr/bin/env bash

printf "%-60s" "Syncing dependencies"
govendor sync
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

sudo setcap 'cap_net_bind_service=+ep' /CodeCollaborate/Server
if [ $? -eq 0 ]; then
    printf "\t%-10s\n" "SETCAP OK"
else
    printf "\t%-10s\n" "SETCAP FAIL"
fi

printf "%-60s" "Restarting systemctl service"
sudo systemctl restart CodeCollaborate.service
if [ $? -eq 0 ]; then
    printf "%-10s\n" "OK"
else
    printf "%-10s\n" "FAIL"
fi
