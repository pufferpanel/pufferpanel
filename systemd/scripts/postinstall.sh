#!/usr/bin/env sh

systemctl daemon-reload

if [ ! -f "/var/lib/pufferpanel/database.db" ]; then
  touch /var/lib/pufferpanel/database.db
fi

pufferpanel --config=/etc/pufferpanel/config.json dbmigrate
exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode

chown -R pufferpanel:pufferpanel /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel /var/www/pufferpanel
