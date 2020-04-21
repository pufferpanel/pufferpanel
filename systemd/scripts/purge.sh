#!/usr/bin/env sh

rm -rf /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel || true
userdel -r  pufferpanel || true
groupdel pufferpanel || true
