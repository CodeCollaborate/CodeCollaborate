#!/usr/bin/env bash
sudo cp CodeCollaborate /etc/init.d/CodeCollaborate
sudo chown root /etc/init.d/CodeCollaborate
sudo chmod 0755 /etc/init.d/CodeCollaborate
sudo systemctl daemon-reload