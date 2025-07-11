#!/usr/bin/env sh

systemctl daemon-reload

mkdir -p /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel /var/www/pufferpanel
if [ ! -f "/var/lib/pufferpanel/database.db" ]; then
  touch /var/lib/pufferpanel/database.db
fi

pufferpanel --config=/etc/pufferpanel/config.json db upgrade
exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode

chown -R pufferpanel:pufferpanel /etc/pufferpanel /var/log/pufferpanel /var/lib/pufferpanel /var/www/pufferpanel

if command -v apparmor_parser >/dev/null 2>&1
then
    apparmor_parser -r /etc/apparmor.d/pufferpanel
fi

chmod o-rx /etc/pufferpanel
chmod o-rx /var/lib/pufferpanel
