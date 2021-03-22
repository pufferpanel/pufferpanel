#!/usr/bin/env sh

rm -rf /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel /var/www/pufferpanel

userdel -r  pufferpanel
exitCode=$?
[ $exitCode -ne 0 && $exitCode -ne 6 ] || exit $exitCode

groupdel pufferpanel
exitCode=$?
[ $exitCode -ne 0 && $exitCode -ne 6 ] || exit $exitCode