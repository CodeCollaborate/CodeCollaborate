#!/usr/bin/env bash
sudo ln -s "$(pwd)/CodeCollaborate.service" /etc/systemd/system/CodeCollaborate.service
sudo systemctl daemon-reload
