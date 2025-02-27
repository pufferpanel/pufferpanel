#!/bin/sh -e

# Source debconf library.
. /usr/share/debconf/confmodule
db_version 2.0

db_title PufferPanel Installation

db_input medium pufferpanel/install_type || true
db_go

db_get pufferpanel/install_type
if [ "$RET" = "panel" ]; then
  db_input medium pufferpanel/create_user || true
  db_go

  db_get pufferpanel/create_user
  if [ "$RET" = "true" ]; then
      until false; do
        db_input medium pufferpanel/email || true
        db_input medium pufferpanel/username || true
        db_input medium pufferpanel/password || true
        db_input medium pufferpanel/password_confirm
        if [ $? -eq 30 ]; then
          break
        fi
        db_go

        db_get pufferpanel/password
        expectedPw=$RET
        db_get pufferpanel/password_confirm
        confirmedPw=$RET

        if [ "$expectedPw" = "$confirmedPw" ]; then
          db_reset pufferpanel/password_confirm
          break
        fi
        db_reset db_reset pufferpanel/password
        db_input medium pufferpanel/password_mismatch
        if [ $? -eq 30 ]; then
          break
        fi
        db_go
      done
  fi
fi

