#!/usr/bin/env sh

useradd --system --home-dir /var/lib/pufferpanel --create-home --user-group pufferpanel

exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode
