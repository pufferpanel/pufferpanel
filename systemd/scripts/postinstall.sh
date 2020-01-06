#!/usr/bin/env sh

systemctl daemon-reload
#useradd --system --home /var/lib/pufferpanel --user-group pufferpanel >/dev/null 2>&1
chown pufferpanel:pufferpanel /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel
