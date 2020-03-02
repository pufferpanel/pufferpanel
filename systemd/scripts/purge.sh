#!/usr/bin/env sh

rm -rf /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel
userdel -r  pufferpanel
groupdel pufferpanel
