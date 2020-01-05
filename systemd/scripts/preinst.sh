#!/usr/bin/env sh

useradd --system --home /var/lib/pufferpanel --user-group pufferpanel
mkdir -p /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel
chown pufferpanel:pufferpanel /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel
