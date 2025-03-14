#!/usr/bin/env sh

rm -rf /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel /var/www/pufferpanel

userdel -r  pufferpanel
exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 6 ] || exit $exitCode

if [ -e /usr/share/debconf/confmodule ]; then
    # Source debconf library.
    . /usr/share/debconf/confmodule
    # Remove my changes to the db.
    db_purge
fi