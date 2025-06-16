#!/usr/bin/env sh

systemctl daemon-reload

if [ ! -f "/var/lib/pufferpanel/database.db" ]; then
  touch /var/lib/pufferpanel/database.db
fi

systemctl is-active --quiet pufferpanel
wasRunning=$?

systemctl stop pufferpanel

pufferpanel --config=/etc/pufferpanel/config.json db upgrade
exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode

chown -R pufferpanel:pufferpanel /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel /var/www/pufferpanel
exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode

if [ $wasRunning -eq 0 ]; then
  systemctl restart pufferpanel
fi

exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode

if command -v apparmor_parser >/dev/null 2>&1
then
    apparmor_parser -r /etc/apparmor.d/unshare
fi
