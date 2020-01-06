#!/usr/bin/env sh

systemctl daemon-reload
chown pufferpanel:pufferpanel /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel
