#!/usr/bin/env sh

if [ -f "/var/lib/pufferpanel/database.db" ]; then
  cp -f /var/lib/pufferpanel/database.db /var/lib/pufferpanel/database-RESTORE.db
fi
